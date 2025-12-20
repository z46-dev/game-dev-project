package definitions

import "fmt"

var (
	ShipConfigs       map[ShipID]*Ship             = make(map[ShipID]*Ship)
	ProjectileConfigs map[ProjectileID]*Projectile = make(map[ProjectileID]*Projectile)
)

func GetByKey[T any, U comparable](confs map[U]*T, key U) (*T, bool) {
	if conf, ok := confs[key]; ok {
		return conf, true
	}

	return nil, false
}

func MustGetByKey[T any, U comparable](confs map[U]*T, key U) *T {
	if conf, ok := confs[key]; ok {
		return conf
	}

	panic(fmt.Errorf("config not found for key: %v", key))
}

const (
	SHIP_TIGER ShipID = iota
	SHIP_HINDENBURG
)

const (
	PROJECTILE_LIGHT_LASER ProjectileID = iota
	PROJECTILE_HEAVY_LASER
	PROJECTILE_PULSE_EMITTER
	PROJECTILE_LIGHT_MISSILE
)
