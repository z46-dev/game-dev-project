package definitions

import "fmt"

var (
	ShipConfigs  map[ShipID]*Ship   = make(map[ShipID]*Ship)
	PlaneConfigs map[PlaneID]*Plane = make(map[PlaneID]*Plane)
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
	SHIP_COLOSSUS ShipID = iota
	SHIP_ENTERPRISE
	SHIP_CHKALOV
	SHIP_PARSEVAL
	SHIP_AQUILA
	SHIP_KAGA
)

const (
	PLANE_VOUGHT_CORSAIR_MKIV PlaneID = iota
	PLANE_FAIREY_BARRACUDA_MKV
	PLANE_F6F_HELLCAT
	PLANE_TBF_AVENGER
	PLANE_SB2C_HELLDIVER
	PLANE_POLIKARPOV_VIT2
	PLANE_SUKHOI_SU2
	PLANE_BF_110C
	PLANE_BF_109G
	PLANE_REGGIANE_RE2001
	PLANE_A6M5_ZERO
	PLANE_B6N_TENZAN
	PLANE_D4Y3_SUISEI
)
