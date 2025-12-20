package definitions

var ProjLightLaser *Projectile = NewProjectile(PROJECTILE_LIGHT_LASER, "Light Laser", 15, 600, 5, 60)

var ProjHeavyLaser *Projectile = NewProjectile(PROJECTILE_HEAVY_LASER, "Heavy Laser", 17.5, 1200, 15, 180)

var PulseEmitter *Projectile = NewProjectile(PROJECTILE_PULSE_EMITTER, "Pulse Emitter", 8, 1000, 0, 180).
	SetSineMovement(2, .5).
	SetExplosion(true, 128, 8, false)