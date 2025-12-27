//go:build ignore

//kage:unit pixels

package main

// Host-supplied uniforms
var Time float
var Camera vec3
var ScreenSize vec2

var StarCenter vec2
var StarRadius float
var StarIntensity float
var StarColor vec3
var StarPulse float
var StarDetail float

// ------------------------- Noise helpers -------------------------

func hash(p vec2) float {
	var h float = dot(p, vec2(127.1, 311.7))
	return fract(sin(h) * 43758.5453)
}

func noise(p vec2) float {
	var i vec2 = floor(p)
	var f vec2 = fract(p)

	var a float = hash(i)
	var b float = hash(i + vec2(1.0, 0.0))
	var c float = hash(i + vec2(0.0, 1.0))
	var d float = hash(i + vec2(1.0, 1.0))

	var u vec2 = f * f * (3.0 - 2.0*f)

	return mix(mix(a, b, u.x), mix(c, d, u.x), u.y)
}

func fbm(p vec2) float {
	var v float = 0.0
	var a float = 0.55
	for i := 0; i < 4; i++ {
		v += a * noise(p)
		p *= 2.2
		a *= 0.5
	}
	return v
}

// Gaussian mask in angle space, used for flare sectors
func angMask(theta, center, width float) float {
	var d float = abs(atan2(sin(theta-center), cos(theta-center)))
	return exp(-(d * d) / (width * width))
}

// ----------------------------- Fragment -----------------------------

func Fragment(pos vec4, _ vec2, _ vec4) vec4 {
	// Position in star-local pixel space
	var toStar vec2 = pos.xy - StarCenter

	// Normalize by StarRadius so rNorm ~ 1 at the photosphere
	var p vec2 = toStar / max(StarRadius, 1.0)
	var rNorm float = length(p)         // radius in "star radii"
	var r2 float = dot(p, p)
	var theta float = atan2(p.y, p.x)

	// Early out: fully transparent far away
	if rNorm > 2.2 {
		return vec4(0.0, 0.0, 0.0, 0.0)
	}

	// Fake 3D sphere normal (for rNorm <= 1)
	var nz float = 0.0
	if rNorm <= 1.0 {
		nz = sqrt(max(0.0, 1.0-r2))
	}
	var n vec3 = normalize(vec3(p.x, p.y, nz))

	// View and light directions
	var view vec3 = normalize(Camera)
	var lightDir vec3 = normalize(vec3(0.4, 0.5, 0.8))

	var ndotl float = max(dot(n, lightDir), 0.0)
	var ndotv float = max(dot(n, view), 0.0)

	// Rim term: bright near edge (as in stars)
	var rim float = pow(1.0-ndotv, 2.0)

	// ---------------------- Photosphere shading ----------------------

	// Rotate surface coords slowly to avoid obvious patterns
	var ca float = cos(Time * 0.05)
	var sa float = sin(Time * 0.05)
	var surfPos vec2 = vec2(
		p.x*ca-p.y*sa,
		p.x*sa+p.y*ca,
	) * StarDetail

	// Drift over time
	surfPos += vec2(Time*0.12, Time*0.08)

	var n1 float = fbm(surfPos * 1.3)
	var n2 float = fbm(surfPos*2.7 + vec2(2.0, -1.5))
	var n3 float = fbm(surfPos*0.8 + vec2(-Time*0.03, Time*0.04))

	var surfaceNoise float = mix(n1, n2, 0.5)
	surfaceNoise = mix(surfaceNoise, n3, 0.4)

	// Base color mix: bright blue-white core, slightly deeper edge
	var coreColor vec3 = vec3(0.96, 0.99, 1.0)
	var midColor vec3 = StarColor
	var edgeColor vec3 = mix(StarColor, vec3(0.25, 0.45, 0.9), 0.5)

	var photosphereMix float = clamp(0.45+0.35*ndotl+0.25*surfaceNoise, 0.0, 1.0)
	var photosphereColor vec3 = mix(edgeColor, coreColor, photosphereMix)
	photosphereColor = mix(photosphereColor, midColor, 0.25)

	// Slight pulsation in brightness
	var corePulse float = 1.0 + StarPulse*0.02*sin(Time*0.5)
	var surfaceBrightness float = (0.8 + 0.6*surfaceNoise) * corePulse
	var rimBoost float = rim * (0.25 + 0.2*surfaceNoise)

	// Photosphere alpha: solid disk with soft limb
	var starAlpha float = smoothstep(1.03, 0.98, rNorm)

	// --------------------------- Corona ---------------------------

	var d float = max(rNorm-1.0, 0.0)
	var coronaPulse float = 1.0 + StarPulse*0.04*sin(Time*0.9)
	var coronaFalloff float = exp(-2.3*d*d)
	var coronaMask float = smoothstep(1.8, 1.0, rNorm)

	var coronaBase float = coronaFalloff * coronaMask * coronaPulse
	var coronaColor vec3 = mix(StarColor, vec3(0.6, 0.85, 1.0), 0.4)

	var coronaRGB vec3 = coronaColor * (coronaBase * 0.3)
	var coronaAlpha float = coronaBase * 0.3

	// ----------------------- Solar flare loops -----------------------

	// Treat star as 3D sphere but build flares in polar space so they
	// wrap around the limb and look like arcs.

	// Three rotating flare sectors
	var pi float = 3.14159265
	var a0 float = Time * 0.10
	var a1 float = 2.0*pi/3.0 + Time*0.08
	var a2 float = 4.0*pi/3.0 + Time*0.06

	var width float = 0.40
	var angM float = angMask(theta, a0, width)
	angM += angMask(theta, a1, width)
	angM += angMask(theta, a2, width)
	angM = pow(clamp(angM, 0.0, 1.0), 1.2)

	// Radial band hugging the limb (slightly inside + slightly outside)
	var flareInner float = smoothstep(0.92, 1.0, rNorm)
	var flareOuter float = 1.0 - smoothstep(1.04, 1.20, rNorm)
	var flareBand float = flareInner * flareOuter
	flareBand *= flareBand

	// Polar coordinates for filaments:
	// x = angle around star, y = distance from surface
	var polar vec2
	polar.x = theta/(2.0*pi) + 0.5           // wrap into [0,1]
	polar.y = (rNorm-1.0) * 4.0             // narrow band around 0

	// Animate arcs: drift in both angle and radial offset
	polar += vec2(Time*0.04, Time*0.09)

	var fN float = fbm(polar * (StarDetail * 0.7))
	var ridge float = 1.0 - abs(2.0*fN-1.0)
	var filament float = pow(ridge, 3.0)    // thin-ish filaments

	var fN2 float = fbm(polar*1.8 + vec2(1.7, -2.1))
	filament *= 0.7 + 0.3*fN2

	// Loop-like pulsation along arcs (no discrete "puffs")
	var loopPulse float = 1.0 + StarPulse*0.08*sin(Time*0.9 + theta*1.4)
	var flareMask float = filament * flareBand * angM * loopPulse
	flareMask = clamp(flareMask, 0.0, 0.9)

	// Separate outside vs inside contributions to hint at 3D loops
	var outsideMask float = smoothstep(1.0, 1.02, rNorm)
	var flareOutside float = flareMask * outsideMask

	// Inside: emphasize near limb / side
	var side float = 1.0 - ndotv
	var surfaceBand float = smoothstep(0.9, 1.0, rNorm)
	var flareInside float = flareMask * surfaceBand * side

	var flareColor vec3 = mix(StarColor, vec3(0.97, 0.99, 1.0), 0.5)

	coronaRGB += flareColor * flareOutside * 0.9
	coronaAlpha += flareOutside * 0.45

	surfaceBrightness += flareInside * 0.8
	rimBoost += flareInside * 0.35

	// --------------------------- Final composite ---------------------------

	var starRGB vec3 = photosphereColor * (surfaceBrightness + rimBoost)

	var rgb vec3 = (starRGB*starAlpha + coronaRGB*coronaAlpha) * StarIntensity
	rgb = clamp(rgb, vec3(0.0, 0.0, 0.0), vec3(1.0, 1.0, 1.0))

	var alpha float = clamp(starAlpha + coronaAlpha, 0.0, 1.0)

	return vec4(rgb, alpha)
}
