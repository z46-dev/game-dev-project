//go:build ignore

//kage:unit pixels

package main

func Fragment(dstPos vec4, srcPos vec2) vec4 {
	// Base sprite with alpha (turret shape)
	base := imageSrc0At(srcPos)

	// If the base image is fully transparent, draw nothing here.
	if base.a <= 0.0 {
		// Returning vec4(0) means fully transparent.
		return vec4(0)
	}

	// Texture image; typically a tileable metal/turret texture.
	tex := imageSrc1At(srcPos)

	// Darken near the shape edge by sampling alpha neighbors.
	a0 := base.a
	aL := imageSrc0At(srcPos + vec2(-1.0, 0.0)).a
	aR := imageSrc0At(srcPos + vec2(1.0, 0.0)).a
	aU := imageSrc0At(srcPos + vec2(0.0, -1.0)).a
	aD := imageSrc0At(srcPos + vec2(0.0, 1.0)).a
	edge := a0 - min(min(aL, aR), min(aU, aD))
	edge = smoothstep(0.0, 0.4, edge)
	darken := mix(vec3(1.0, 1.0, 1.0), vec3(0.7, 0.7, 0.7), edge)

	// Apply the texture color but keep the base sprite's alpha (mask).
	tex.rgb *= darken
	return vec4(tex.rgb, base.a)
}
