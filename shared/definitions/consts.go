package definitions

const (
	EntityTypeShip EntityType = iota
	EntityTypeProjectile
	EntityTypeAsteroid
	EntityTypeStation
)

const (
	ShipClassificationFighter   ShipClassification = iota // Small, fast, and agile ships designed for close-range support and dogfighting
	ShipClassificationBomber                              // Bigger and slower than fighters, but still small. Designed to deliver heavy payloads to larger targets
	ShipClassificationCorvette                            // Small capital ships designed for patrol and escort duties
	ShipClassificationFrigate                             // Medium-sized capital ships designed for multi-role combat and support
	ShipClassificationDestroyer                           // Large capital ships designed for frontline combat and heavy firepower
	ShipClassificationCarrier                             // Very large capital ships designed to deploy and support smaller craft
)
