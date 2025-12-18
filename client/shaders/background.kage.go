//go:build ignore

//kage:unit pixels

package main

var Time float      // Time in ticks
var Camera vec3     // (x, y, zoom)
var ScreenSize vec2 // (w, h)

const Radius float = 100
const PI = 3.141592653589793
const TAU = 6.283185307179586
const RADIANS_120 = 2.0943951023931953
const RADIANS_240 = 4.1887902047863905

func hash(p vec2) float {
	return fract(sin(dot(p, vec2(127.1, 311.7))) * 32938.5453)
}

func noise(p vec2) float {
	i := floor(p)
	f := fract(p)
	u := f * f * (vec2(3.0, 3.0) - 2.0*f)

	return mix(
		mix(hash(i+vec2(0.0, 0.0)), hash(i+vec2(1.0, 0.0)), u.x),
		mix(hash(i+vec2(0.0, 1.0)), hash(i+vec2(1.0, 1.0)), u.x),
		u.y,
	)
}

func fbmWithTime(p vec2, t float) float {
	var v, a float = 0.0, 0.5
	for i := 0; i < 4; i++ {
		v += a * noise(p+vec2(t*0.02, t*0.04))
		p = p*2.0 + vec2(5.2, 1.3)
		a *= 0.5
	}
	return v
}

func Fragment(dstPos vec4, srcPos vec2, _ vec4) vec4 {
    var world vec2 = (dstPos.xy-ScreenSize*0.5)/Camera.z + Camera.xy

    var (
        rel vec2 = (world) / Radius
        d float = length(rel)
    )

    if d < 1 {
        var hue float = (((atan2(rel.y, rel.x) + Time * 0.001) / TAU) + 0.5) * TAU
        return vec4(0.5 + 0.5 * cos(hue), 0.5 + 0.5 * cos(hue + RADIANS_120), 0.5 + 0.5 * cos(hue + RADIANS_240), 1.0)
    } else {
        var (
            n float = fbmWithTime(world*0.1, Time)
            c vec3 = vec3(n)
        )
        return vec4(c, 1.0)
    }
}