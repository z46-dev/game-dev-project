package definitions

var ProjLightLaser *Projectile = NewProjectile(PROJECTILE_LIGHT_LASER, "Light Laser", 20, 3000, 5, 15)

var ProjHeavyLaser *Projectile = NewProjectile(PROJECTILE_HEAVY_LASER, "Heavy Laser", 20, 4500, 15, 180)

var PulseEmitter *Projectile = NewProjectile(PROJECTILE_PULSE_EMITTER, "Pulse Emitter", 15, 1000, 0, 180).
	SetSineMovement(2, .5).
	SetExplosion(true, 128, 8, false)