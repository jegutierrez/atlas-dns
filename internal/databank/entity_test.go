package databank_test

import (
	"encoding/json"
	"testing"

	"github.com/jegutierrez/atlas-dns/internal/databank"
)

func TestRobotStatusUnmarshal(t *testing.T) {
	tt := []struct {
		description         string
		requestBody         string
		expectedRobotStatus databank.RobotState
		expectedErr         error
	}{
		{
			description: "valid robot status request",
			requestBody: `{
			"x": "123.12",
			"y": "456.56",
			"z": "789.89",
			"vel": "20.0"
			}`,
			expectedRobotStatus: databank.RobotState{
				X:        123.12,
				Y:        456.56,
				Z:        789.89,
				Velocity: 20.0,
			},
			expectedErr: nil,
		},
		{
			description: "invalid robot status request",
			requestBody: `"invalid_body"`,
			expectedErr: databank.ErrInvalidRobotStatus,
		},
		{
			description: "invalid robot status request, missing field",
			requestBody: `{
			"y": "456.56",
			"z": "789.89",
			"vel": "20.0"
			}`,
			expectedErr: databank.ErrInvalidRobotStatus,
		},
		{
			description: "invalid robot X coordinate request",
			requestBody: `{
			"x": "XXX.XX",
			"y": "456.56",
			"z": "789.89",
			"vel": "20.0"
			}`,
			expectedErr: databank.ErrInvalidCoordinate,
		},
		{
			description: "invalid robot Y coordinate request",
			requestBody: `{
			"x": "123.2",
			"y": "XXX.X",
			"z": "789.89",
			"vel": "20.0"
			}`,
			expectedErr: databank.ErrInvalidCoordinate,
		},
		{
			description: "invalid robot Z coordinate request",
			requestBody: `{
			"x": "123.12",
			"y": "456.56",
			"z": "XXX",
			"vel": "20.0"
			}`,
			expectedErr: databank.ErrInvalidCoordinate,
		},
		{
			description: "invalid robot velocity request",
			requestBody: `{
			"x": "123.12",
			"y": "456.56",
			"z": "789.89",
			"vel": "XXXX.X"
			}`,
			expectedErr: databank.ErrInvalidVelocity,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			var unmarshalResult databank.RobotState

			err := json.Unmarshal([]byte(tc.requestBody), &unmarshalResult)
			if err == nil && tc.expectedErr != nil {
				t.Fatalf("an error was expected: %s and got nil", tc.expectedErr.Error())
			}
			if err != nil {
				if tc.expectedErr == nil {
					t.Fatalf("unexpected error return: %s", err.Error())
				} else if err.Error() != tc.expectedErr.Error() {
					t.Fatalf("unexpected error return, want: %s, got: %s", tc.expectedErr.Error(), err.Error())
				}
			}
			if err == nil && unmarshalResult.String() != tc.expectedRobotStatus.String() {
				t.Fatalf("unexpected result from robot status unmarshal, want: %s and got: %s", tc.expectedRobotStatus, unmarshalResult)
			}
		})
	}
}
