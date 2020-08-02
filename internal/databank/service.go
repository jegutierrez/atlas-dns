package databank

// Service abstracts the databank's location service.
type Service interface {
	GetNearestDatabankLocation(robot RobotState, sectorID int) (Location, error)
}

// NewLocationService builds a databank's location service.
func NewLocationService() Service {
	return &LocatorService{}
}

// LocatorService holds the business logic to calculate databank locations.
type LocatorService struct{}

// GetNearestDatabankLocation computes the nearest databank from the given coordinates and velocity.
func (l *LocatorService) GetNearestDatabankLocation(robot RobotState, sectorID int) (Location, error) {
	fSectorID := float64(sectorID)
	loc := robot.X*fSectorID + robot.Y*fSectorID + robot.Z*fSectorID + robot.Velocity
	return Location{Location: loc}, nil
}
