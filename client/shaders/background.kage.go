//go:build ignore

//kage:unit pixels

package main

var Time float      // Time in ticks
var Camera vec3     // (x, y, zoom)
var ScreenSize vec2 // (w, h)

const (
	brownianScale  = 0.5825
	detailLevel    = 2.3
	colorScaleBase = 0.25
	colorScaleMix  = 0.5
	waveStrength   = 0.1
	worldScale     = 0.002
)

func random(point vec2) float {
	var additive vec2 = vec2(12.9898, 78.233)
	return fract(sin(dot(point, additive)) * 43758.5453123)
}

func noise(point vec2) float {
	var floorPoint vec2 = floor(point)
	var fraction vec2 = fract(point)
	var smoothedFraction vec2 = fraction*fraction*(vec2(3.0, 3.0)-2.0*fraction)

	var a float = random(floorPoint)
	var b float = random(floorPoint + vec2(1.0, 0.0))
	var c float = random(floorPoint + vec2(0.0, 1.0))
	var d float = random(floorPoint + vec2(1.0, 1.0))

	return mix(a, b, smoothedFraction.x) +
		(c-a)*smoothedFraction.y*(1.0-smoothedFraction.x) +
		(d-b)*smoothedFraction.x*smoothedFraction.y
}

func fbm(point vec2) float {
	var vertexValue float = 0.0
	var fragmentAlpha float = brownianScale

	var cos0_5 float = cos(0.5)
	var sin0_5 float = sin(0.5)
	var rot mat2 = mat2(cos0_5, sin0_5, -sin0_5, cos0_5)
	var shift vec2 = vec2(100.0, 100.0)

	for i := 0; i < 5; i++ {
		vertexValue += fragmentAlpha * noise(point)
		point = rot*point*detailLevel + shift
		fragmentAlpha *= brownianScale
	}

	return colorScaleBase + colorScaleMix*vertexValue
}

func Fragment(dstPos vec4, srcPos vec2, _ vec4) vec4 {
	var world vec2 = (dstPos.xy-ScreenSize*0.5)/Camera.z + Camera.xy
	var point vec2 = world * worldScale
	var time float = Time / 1024.0

	var color vec3 = vec3(0.22, 0.5, 0.67)

	var position vec2
	position.x = fbm(point)
	position.y = fbm(point + vec2(1.0, 1.0))

	var offset vec2
	offset.x = fbm(point + position + vec2(1.7, 9.2) + vec2(time, time))
	offset.y = fbm(point + position + vec2(8.3, 2.8) + vec2(time*0.7, time*0.7))

	offset.x += fbm(point+offset+vec2(0.0, 1.0)+vec2(time, time)) * waveStrength
	offset.y += fbm(point+offset+vec2(1.0, 0.0)+vec2(time, time)) * waveStrength

	var f float = fbm(point + offset)

	var deepWater vec3 = vec3(0.0, 0.4, 0.6)
	var shallowWater vec3 = vec3(0.0, 0.6, 0.85)
	var foam vec3 = vec3(1.0, 1.0, 1.0)

	color = mix(deepWater, shallowWater, clamp((f*f)*4.0, 0.0, 1.0))
	color = mix(color, shallowWater, clamp(abs(offset.x), 0.0, 1.0)*0.6)
	color = mix(color, foam, clamp(length(position), 0.0, 1.0)*0.3)

	var colorScaleSum float = colorScaleBase + colorScaleMix
	if colorScaleSum > 0.6 {
		var foamValue float = fbm(point+offset) * (0.5 + 0.5*sin(time*0.25+point.x*point.y))
		if foamValue > colorScaleSum {
			color = mix(color, foam, clamp((foamValue-colorScaleSum)*100.0, 0.0, 1.0))
		}
	}

	var final float = f*f*f + 0.6*f*f + 0.5*f
	return vec4(final*color, 1.0)
}
