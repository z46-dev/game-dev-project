//go:build ignore

//kage:unit pixels

package main

var Time float      // Time in ticks
var Camera vec3     // (x, y, zoom)
var ScreenSize vec2 // (w, h)

func hash(p vec2) float {
	return fract(sin(dot(p, vec2(127.1, 311.7))) * 43758.5453)
}

func noise(p vec2) float {
	i := floor(p)
	f := fract(p)
	u := f*f*(vec2(3.0, 3.0)-2.0*f)

	return mix(
		mix(hash(i+vec2(0.0, 0.0)), hash(i+vec2(1.0, 0.0)), u.x),
		mix(hash(i+vec2(0.0, 1.0)), hash(i+vec2(1.0, 1.0)), u.x),
		u.y,
	)
}

func fbm(p vec2) float {
	var v float = 0.0
	var a float = 0.5
	for i := 0; i < 5; i++ {
		v += a * noise(p)
		p = p*2.0 + vec2(11.2, 7.4)
		a *= 0.5
	}
	return v
}

func starfield(p vec2) float {
	var grid vec2 = floor(p)
	var h float = hash(grid)
	var sparkle float = smoothstep(0.9975, 1.0, h)
	var twinkle float = 0.75 + 0.25*sin(h*50.0+Time*0.02)
	return sparkle * twinkle
}

func swirl(p vec2) vec3 {
	var angle float = atan2(p.y, p.x)
	var radius float = length(p)
	var band float = sin(angle*3.0+radius*6.0-Time*0.0006)
	var mask float = smoothstep(0.55, 0.75, band) * smoothstep(1.2, 0.4, radius)
	return vec3(0.12, 0.04, 0.2) * mask
}

func Fragment(dstPos vec4, srcPos vec2, _ vec4) vec4 {
	var world vec2 = (dstPos.xy-ScreenSize*0.5)/Camera.z + Camera.xy
	var uv vec2 = world * 0.002

	var base vec3 = vec3(0.02, 0.03, 0.06)
	var nebula1 float = fbm(uv*2.0 + vec2(Time*0.0003, -Time*0.0002))
	var nebula2 float = fbm(uv*3.5 + vec2(-Time*0.0002, Time*0.00035))

	var fog vec3 = vec3(0.08, 0.02, 0.12) * nebula1
	fog += vec3(0.02, 0.08, 0.12) * nebula2
	fog *= 0.8

	var swirlSeed vec2 = floor(world*0.0008)
	var swirlRand float = hash(swirlSeed)
	var swirlMask float = smoothstep(0.985, 1.0, swirlRand)
	var swirlColor vec3 = swirl(fract(world*0.0008)*2.0-vec2(1.0, 1.0)) * swirlMask

	var stars float = starfield(world*0.08)
	var stars2 float = starfield(world*0.16) * 0.5
	var starColor vec3 = vec3(0.9, 0.95, 1.0) * (stars + stars2)

	var color vec3 = base + fog + swirlColor + starColor
	return vec4(color, 1.0)
}
