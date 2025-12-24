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
	var p0 vec2 = (position.xy - StarCenter) / max(StarRadius, 1.0)
	var r float = length(p0)
	var r2 float = dot(p0, p0)
	var theta float = atan2(p0.y, p0.x)

	var nx float = p0.x
	var ny float = p0.y
	var nz float = sqrt(max(0.0, 1.0-r2))
	var n vec3 = normalize(vec3(nx, ny, nz))
	var view vec3 = normalize(Camera)
	var ndotl float = max(dot(n, view), 0.0)

	var rim float = smoothstep(0.8, 1.0, r)
	var limb float = mix(0.65, 1.0, ndotl)

	var ca float = cos(Time * 0.05)
	var sa float = sin(Time * 0.05)
	var q vec2 = vec2(p0.x*ca-p0.y*sa, p0.x*sa+p0.y*ca) * StarDetail
	var drift float = 0.2 + 0.1*sin(Time*0.12)
	q += vec2(Time*0.1, Time*0.07) + vec2(-p0.y, p0.x)*drift

	var n1 float = fbm(q*1.5)
	var n2 float = fbm(q*2.6 + vec2(1.7, -2.3))
	var n3 float = fbm(q*0.7 + vec2(-Time*0.05, Time*0.03))
	var n4 float = fbm(q*4.0 + vec2(Time*0.18, -Time*0.14))
	var surfaceNoise float = mix(n1, n2, 0.5)
	surfaceNoise = mix(surfaceNoise, n3, 0.35)

	var baseMix float = clamp(0.5+0.4*ndotl+0.2*surfaceNoise, 0.0, 1.0)
	var coreColor vec3 = vec3(0.95, 0.98, 1.0)
	var edgeColor vec3 = mix(StarColor, vec3(0.2, 0.4, 0.9), 0.5)
	var starColor vec3 = mix(edgeColor, coreColor, baseMix)
	starColor = mix(starColor, StarColor, 0.25+0.25*surfaceNoise)
	starColor = mix(starColor, edgeColor, 0.18*(1.0-n3))

	var starAlpha float = smoothstep(1.02, 0.98, r)
	var corePulse float = 1.0 + StarPulse*0.015*sin(Time*0.35+surfaceNoise*1.5)
	var spot float = mix(0.75, 1.2, n3) * mix(0.9, 1.1, n4)
	var surface float = (0.7 + 0.6*surfaceNoise) * limb * corePulse * spot
	var rimBoost float = rim * (0.2 + 0.2*surfaceNoise)

	var coronaPulse float = 1.0 + StarPulse*0.03*sin(Time*0.6)
	var rr float = r
	var coronaFalloff float = exp(-2.2*(rr-1.0)*(rr-1.0))
	var coronaMask float = smoothstep(1.6, 1.0, r)
	var corona float = coronaFalloff * coronaMask * coronaPulse
	var coronaColor vec3 = mix(StarColor, vec3(0.5, 0.8, 1.0), 0.35)
	var coronaBase float = corona * 0.35
	var coronaAlpha float = coronaBase * 0.35
	var coronaRGB vec3 = coronaColor * coronaBase

	var rot float = Time * 0.02
	var angVec vec2 = vec2(cos(theta+rot*0.6), sin(theta+rot*0.6))
	var angNoise float = fbm(angVec*2.0 + vec2(Time*0.03, -Time*0.02))
	var angGateA float = smoothstep(0.6, 0.85, fbm(angVec*1.6 + vec2(Time*0.07, -Time*0.05)))
	var angGateB float = smoothstep(0.55, 0.8, fbm(angVec*2.8 + vec2(Time*0.04, Time*0.03)))
	var angMask float = angGateA * angGateB * smoothstep(0.45, 0.75, 1.0-abs(2.0*angNoise-1.0))

	var thetaWarp float = theta + (angNoise-0.5)*0.7 + (fbm(angVec*3.1 + vec2(Time*0.02, Time*0.025))-0.5)*0.35
	var rWarp float = r + (fbm(vec2(thetaWarp*1.7, Time*0.06))-0.5)*0.1 + sin(thetaWarp*2.4+Time*0.18)*0.05

	var coord vec2 = vec2(thetaWarp*StarDetail*1.8, (rWarp-1.0)*StarDetail*6.0)
	coord += vec2(Time*0.18, -Time*0.12)
	var nF float = fbm(coord*1.0)
	var ridge float = 1.0 - abs(2.0*nF-1.0)
	var filament float = pow(ridge, 4.2)
	var fine float = fbm(coord*1.9 + vec2(1.7, -2.1))
	var filamentMask float = filament * (0.55 + 0.45*fine)

	var innerMask float = smoothstep(0.8, 1.0, rWarp)
	var outerMask float = 1.0 - smoothstep(1.3, 1.6, rWarp)
	var radialMask float = innerMask * outerMask
	var innerSurface float = 1.0 - smoothstep(0.7, 1.0, rWarp)
	var flarePulse float = 1.0 + StarPulse*0.08*sin(Time*1.1+fine*1.6+angNoise*2.0)
	var flareStrength float = filamentMask * radialMask * angMask * flarePulse
	var flareColor vec3 = mix(StarColor, vec3(0.85, 0.92, 1.0), 0.45)
	var flareInside float = flareStrength * innerSurface * 0.5
	var flareOutside float = flareStrength
	coronaRGB += flareColor * flareOutside * 0.9
	coronaAlpha += flareOutside * 0.25
	rimBoost += flareInside * 0.22

	var pi float = 3.14159265
	var a0 float = 0.0 + Time*0.12
	var a1 float = 2.0*pi/3.0 + Time*0.09
	var a2 float = 4.0*pi/3.0 + Time*0.07
	var flareWidth float = 0.32
	var d0 float = abs(atan2(sin(theta-a0), cos(theta-a0)))
	var d1 float = abs(atan2(sin(theta-a1), cos(theta-a1)))
	var d2 float = abs(atan2(sin(theta-a2), cos(theta-a2)))
	var angMaskArc float = exp(-((d0/flareWidth)*(d0/flareWidth)))
	angMaskArc += exp(-((d1/flareWidth)*(d1/flareWidth)))
	angMaskArc += exp(-((d2/flareWidth)*(d2/flareWidth)))
	angMaskArc = clamp(angMaskArc*1.15, 0.0, 1.0)

	var flareInner float = smoothstep(0.85, 1.0, r)
	var flareOuter float = 1.0 - smoothstep(1.3, 1.7, r)
	var flareRadial float = flareInner * flareOuter
	var surfaceMask float = 1.0 - smoothstep(0.7, 1.0, r)

	var tangent vec2 = normalize(vec2(-p0.y, p0.x) + vec2(0.0001, 0.0))
	var flareCoord vec2 = tangent * (StarDetail * 2.0) + p0 * 0.4
	flareCoord += vec2(Time*0.25, Time*0.18)
	var f0 float = fbm(flareCoord*1.2)
	var ridgeF float = 1.0 - abs(2.0*f0-1.0)
	var filamentF float = pow(ridgeF, 4.0)
	var f1 float = fbm(flareCoord*2.1 + vec2(1.5, -2.2))
	filamentF *= 0.6 + 0.4*f1

	var flareBase float = filamentF * flareRadial * angMaskArc
	flareBase = clamp(flareBase*1.6 + 0.08*flareRadial*angMaskArc, 0.0, 1.0)

	var phase float = dot(flareCoord, vec2(0.7, 1.3)) * 1.4 + Time*0.8
	var burstPhase float = fract(phase)
	var burst float = smoothstep(0.0, 0.18, burstPhase) * (1.0 - smoothstep(0.35, 0.55, burstPhase))
	var burstGate float = smoothstep(0.35, 0.7, fbm(flareCoord*0.6 + vec2(Time*0.08, -Time*0.06)))
	burst *= burstGate

	var arcPulse float = 1.0 + StarPulse*0.1*sin(Time*1.1)
	var flareArc float = flareBase * arcPulse
	var flareBurst float = flareBase * burst * 2.0
	var flareTotal float = clamp(flareArc + flareBurst, 0.0, 1.0)

	var flareColorArc vec3 = mix(StarColor, vec3(0.95, 0.98, 1.0), 0.5)
	var flareOutsideMask float = smoothstep(0.98, 1.05, r)
	var flareOutsideArc float = flareTotal * flareOutsideMask
	var flareInsideArc float = flareTotal * surfaceMask * 0.5
	coronaRGB += flareColorArc * flareOutsideArc * 0.8
	coronaAlpha += flareOutsideArc * 0.4
	surface += flareInsideArc * 0.35
	rimBoost += flareInsideArc * 0.1

	var starRGB vec3 = starColor * (surface + rimBoost)
	var rgb vec3 = (starRGB*starAlpha + coronaRGB*coronaAlpha) * StarIntensity
	rgb = clamp(rgb, vec3(0.0, 0.0, 0.0), vec3(1.0, 1.0, 1.0))
	var alpha float = clamp(starAlpha + coronaAlpha, 0.0, 1.0)
	return vec4(rgb, alpha)
}
