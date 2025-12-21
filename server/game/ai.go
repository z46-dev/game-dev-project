package game

type (
	Candidate struct {
		Ship *Ship
		Dist float64
	}

	ShipAI struct {
		Ship   *Ship
		Target *Ship
	}
)

func ClosestShip(ships []*Candidate) (closest *Candidate) {
	if len(ships) == 0 {
		return nil
	}

	closest = ships[0]
	for _, ship := range ships {
		if ship.Dist < closest.Dist {
			closest = ship
		}
	}

	return
}

func SelectShipsAroundMe(me *Ship, closedSearchRange float64) []*Candidate {
	var found []*Candidate

	return found
}
