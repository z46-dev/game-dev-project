package definitions

const (
	EntityTypeShip EntityType = iota
	EntityTypeProjectile
	EntityTypeAsteroid
	EntityTypeStation
)

const (
	MovementTypeLinear   ProjectileMovementPattern = iota // Default straight-line movement
	MovementTypeSineWave                                  // Sine wave movement (factor of amplitude and frequency)
	MovementTypeHoming                                    // Homing movement towards a target (factor of turn rate)
)

const (
	HardpointDrawLayerAboveHull HardpointDrawLayer = iota // Draw the hardpoint above the ship hull
	HardpointDrawLayerBelowHull                           // Draw the hardpoint below the ship hull
	HardpointDrawLayerHidden                              // Do not draw the hardpoint
)

const (
	ShipClassificationFighter   ShipClassification = iota // Small, fast, and agile ships designed for close-range support and dogfighting
	ShipClassificationBomber                              // Bigger and slower than fighters, but still small. Designed to deliver heavy payloads to larger targets
	ShipClassificationCorvette                            // Small capital ships designed for patrol and escort duties
	ShipClassificationFrigate                             // Medium-sized capital ships designed for multi-role combat and support
	ShipClassificationDestroyer                           // Large capital ships designed for frontline combat and heavy firepower
)
