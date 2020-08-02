package databank_test

import (
	"fmt"
	"testing"

	"github.com/jegutierrez/atlas-dns/internal/databank"
)

func TestDatabankLocationCalculator(t *testing.T) {
	tt := []struct {
		description              string
		robotState               databank.RobotState
		sectorID                 int
		expectedLocationResponse float64
		expectedErr              error
	}{
		{
			description: "valid robot location request, expected 32.0",
			robotState: databank.RobotState{
				X:        2.0,
				Y:        2.0,
				Z:        2.0,
				Velocity: 20.0,
			},
			sectorID:                 2,
			expectedLocationResponse: 32.0,
			expectedErr:              nil,
		},
		{
			description: "valid robot location request, expected ",
			robotState: databank.RobotState{
				X:        12.12,
				Y:        43.78,
				Z:        121.0,
				Velocity: 20.0,
			},
			sectorID:                 5,
			expectedLocationResponse: 904.5,
			expectedErr:              nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			service := databank.NewLocationService()
			databank, err := service.GetNearestDatabankLocation(tc.robotState, tc.sectorID)
			if err == nil && tc.expectedErr != nil {
				t.Fatalf("an error was expected: %s and got nil", tc.expectedErr.Error())
			}
			if err == nil && fmt.Sprintf("%f", databank.Location) != fmt.Sprintf("%f", tc.expectedLocationResponse) {
				t.Fatalf("unexpected databank location, want: %f and got: %f", tc.expectedLocationResponse, databank.Location)
			}
		})
	}
}
