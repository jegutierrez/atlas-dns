package databank

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	// ErrInvalidRobotStatus returned when the request from a robot can't be unmarshalled into a struct.
	ErrInvalidRobotStatus = errors.New("Invalid robot status, X,Y,Z coordinates and Velocity are required")
	// ErrInvalidCoordinate returned when a given coordinate can't be transformed into a float64.
	ErrInvalidCoordinate = errors.New("Invalid coordinate, a floating point number was expected")
	// ErrInvalidVelocity returned when a given velocity value can't be transformed into a float64.
	ErrInvalidVelocity = errors.New("Invalid velocity, a floating point number was expected")
)

// RobotState represents a request from a robot.
type RobotState struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z"`
	Velocity float64 `json:"vel"`
}

// UnmarshalJSON transform bytes int a RobotState struct.
func (r *RobotState) UnmarshalJSON(data []byte) error {
	type robotStateRequest struct {
		X   string `json:"x"`
		Y   string `json:"y"`
		Z   string `json:"z"`
		Vel string `json:"vel"`
	}
	var request robotStateRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return ErrInvalidRobotStatus
	}
	if request.X == "" || request.Y == "" || request.Z == "" || request.Vel == "" {
		return ErrInvalidRobotStatus
	}

	var x, y, z, vel float64
	var err error

	if x, err = strconv.ParseFloat(request.X, 64); err != nil {
		return ErrInvalidCoordinate
	}
	if y, err = strconv.ParseFloat(request.Y, 64); err != nil {
		return ErrInvalidCoordinate
	}
	if z, err = strconv.ParseFloat(request.Z, 64); err != nil {
		return ErrInvalidCoordinate
	}
	if vel, err = strconv.ParseFloat(request.Vel, 64); err != nil {
		return ErrInvalidVelocity
	}

	r.X = x
	r.Y = y
	r.Z = z
	r.Velocity = vel

	return nil
}

func (r RobotState) String() string {
	return fmt.Sprintf(`{"x":%f, "y": %f, "z": %f, "vel": %f}`, r.X, r.Y, r.Z, r.Velocity)
}

// Location represents a databank location info.
type Location struct {
	Location float64 `json:"loc"`
}

// LocationAPIError error returned to clients when the databank location fails.
type LocationAPIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func newLocationAPIError(statusCode int, msg string) LocationAPIError {
	return LocationAPIError{
		StatusCode: statusCode,
		Message:    msg,
	}
}
