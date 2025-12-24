//go:build ignore

//kage:unit pixels

package main

var Time float      // Time in ticks
var Camera vec3     // (x, y, zoom)
var ScreenSize vec2 // (w, h)

var StarCenter vec2
var StarRadius float
var StarIntensity float
var StarColor vec3
var StarPulse float
var StarDetail float

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
		p *= 2.1
		a *= 0.5
	}
	return v
}

func Fragment(position vec4, texCoord vec2, srcColor vec4) vec4 {
	var p vec2 = (position.xy - StarCenter) / max(StarRadius, 1.0)

	var r float = length(p)
	var r2 float = dot(p, p)

	var nx float = p.x
	var ny float = p.y
	var nz float = sqrt(max(0.0, 1.0-r2))
	var n vec3 = normalize(vec3(nx, ny, nz))
	var view vec3 = normalize(Camera)
	var ndotl float = max(dot(n, view), 0.0)

	var rim float = smoothstep(0.8, 1.0, r)
	var limb float = mix(0.65, 1.0, ndotl)

	var ca float = cos(Time * 0.05)
	var sa float = sin(Time * 0.05)
	var q vec2 = vec2(p.x*ca-p.y*sa, p.x*sa+p.y*ca) * StarDetail
	q += vec2(Time*0.1, Time*0.07)

	var n1 float = fbm(q*1.5)
	var n2 float = fbm(q*2.6 + vec2(1.7, -2.3))
	var surfaceNoise float = mix(n1, n2, 0.5)

	var baseMix float = clamp(0.5+0.4*ndotl+0.2*surfaceNoise, 0.0, 1.0)
	var coreColor vec3 = vec3(0.95, 0.98, 1.0)
	var edgeColor vec3 = mix(StarColor, vec3(0.2, 0.4, 0.9), 0.5)
	var starColor vec3 = mix(edgeColor, coreColor, baseMix)
	starColor = mix(starColor, StarColor, 0.25+0.25*surfaceNoise)

	var starAlpha float = smoothstep(1.02, 0.98, r)
	var surface float = (0.7 + 0.6*surfaceNoise) * limb
	var rimBoost float = rim * (0.2 + 0.2*surfaceNoise)

	var coronaFalloff float = exp(-2.6*(r-1.0)*(r-1.0))
	var coronaMask float = smoothstep(2.0, 1.0, r)
	var cNoise float = fbm(p*0.5 + vec2(Time*0.02, -Time*0.015))
	var outerMask float = smoothstep(1.0, 1.05, r)
	var glowPulse float = 1.0 + StarPulse*0.06*sin(Time*0.6+cNoise*2.0)
	var corona float = coronaFalloff * coronaMask * outerMask * (0.65 + 0.35*cNoise) * glowPulse
	var coronaColor vec3 = mix(StarColor, vec3(0.5, 0.8, 1.0), 0.35)
	var coronaAlpha float = corona * 0.6
	var coronaRGB vec3 = coronaColor * corona

	var ang float = atan2(p.y, p.x)
	var bend float = fbm(vec2(ang*1.5, r*2.2+Time*0.05))
	var rWarp float = r + (bend-0.5)*0.12
	var filamentNoise float = fbm(vec2(ang*3.5+bend*2.0, rWarp*3.0-Time*0.08))
	var filament float = smoothstep(0.6, 0.85, filamentNoise)
	var filamentRad float = smoothstep(0.95, 1.15, rWarp) * smoothstep(2.0, 1.1, rWarp)
	var filamentStrength float = filament * filamentRad * (0.5 + 0.5*bend)
	var filamentColor vec3 = mix(StarColor, vec3(0.65, 0.85, 1.0), 0.45)
	coronaRGB += filamentColor * filamentStrength * 1.2
	coronaAlpha += filamentStrength * 0.3

	var rimFil float = filament * smoothstep(1.0, 0.9, rWarp)
	rimBoost += rimFil * 0.18

	var starRGB vec3 = starColor * (surface + rimBoost)

	var rgb vec3 = (starRGB*starAlpha + coronaRGB*coronaAlpha) * StarIntensity
	var alpha float = clamp(starAlpha + coronaAlpha, 0.0, 1.0)
	return vec4(rgb, alpha)
}
