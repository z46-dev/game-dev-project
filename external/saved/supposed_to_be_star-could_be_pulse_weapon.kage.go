//go:build ignore

//kage:unit pixels

package main

var Time float      // Time in ticks
var Camera vec3     // (x, y, zoom)
var ScreenSize vec2 // (w, h)
var StarCenter vec2 // World-space position of the star center
var StarRadius float
var StarIntensity float
var StarColor vec3 // Base star color (RGB)
var StarPulse float
var StarDetail float

func hash21(p vec2) float {
	var h float = dot(p, vec2(127.1, 311.7))
	return fract(sin(h) * 43758.5453123)
}

func noise2(p vec2) float {
	var i vec2 = floor(p)
	var f vec2 = fract(p)
	var a float = hash21(i)
	var b float = hash21(i + vec2(1.0, 0.0))
	var c float = hash21(i + vec2(0.0, 1.0))
	var d float = hash21(i + vec2(1.0, 1.0))
	var u vec2 = f * f * (3.0 - 2.0*f)
	return mix(mix(a, b, u.x), mix(c, d, u.x), u.y)
}

func fbm(p vec2) float {
	var v float = 0.0
	var a float = 0.5
	for i := 0; i < 5; i++ {
		v += a * noise2(p)
		p *= 2.05
		a *= 0.5
	}
	return v
}

func starField(world vec2) vec3 {
	var p vec2 = world * 0.02
	var n float = fbm(p*3.5 + vec2(Time*0.01, -Time*0.007))
	var twinkle float = sin(Time*0.08+n*8.0)*0.5 + 0.5
	var base vec3 = vec3(0.02, 0.04, 0.07)
	var stars float = smoothstep(0.86, 1.0, n) * (0.6 + 0.4*twinkle)
	return base + stars*vec3(0.35, 0.45, 0.6)
}

func starCore(world vec2) vec3 {
	var toCenter vec2 = world - StarCenter
	var r float = length(toCenter)
	var dir vec2 = toCenter / (r + 1e-4)
	var t float = Time * 0.02

	var radius float = max(StarRadius, 1.0)
	var inner float = exp(-r*r/(radius*radius*0.4))
	var mid float = exp(-r*r/(radius*radius*1.6))
	var halo float = exp(-r*r/(radius*radius*6.0))

	var swirl float = fbm(dir*vec2(1.3, 0.7)*StarDetail + vec2(t, -t*1.1))
	var band float = sin((dir.x-dir.y)*6.0 + swirl*6.0 + t*4.0)*0.5 + 0.5
	var rays float = pow(abs(dir.x*dir.y), 0.35)
	var rayMask float = smoothstep(0.45, 1.0, rays)

	var pulse float = 0.85 + 0.15*sin(Time*StarPulse*0.1)
	var detail float = mix(0.6, 1.2, swirl) * (0.8 + 0.2*band)

	var coreColor vec3 = mix(StarColor, vec3(0.6, 0.75, 1.0), 0.35)
	var glow vec3 = coreColor * (inner*2.0 + mid*1.1) * detail
	var corona vec3 = vec3(0.2, 0.4, 0.9) * (halo*0.7 + rayMask*halo*1.6)

	return (glow + corona) * StarIntensity * pulse
}

func Fragment(dstPos vec4, srcPos vec2, _ vec4) vec4 {
	var world vec2 = (dstPos.xy-ScreenSize*0.5)/Camera.z + Camera.xy
	var bg vec3 = starField(world)
	var star vec3 = starCore(world)
	var color vec3 = bg + star
	return vec4(color, 1.0)
}
