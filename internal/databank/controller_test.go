package databank_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jegutierrez/atlas-dns/internal/app"
	"github.com/jegutierrez/atlas-dns/internal/databank"
)

type mockLocationService struct {
	response float64
	err      bool
}

func (l mockLocationService) GetNearestDatabankLocation(robot databank.RobotState, sectorID int) (databank.Location, error) {
	if l.err {
		return databank.Location{}, errors.New("some error")
	}
	return databank.Location{Location: l.response}, nil
}

func buildMockLocationService(response float64, err bool) databank.Service {
	return mockLocationService{
		response: response,
		err:      err,
	}
}

func TestLocationHttpHandler(t *testing.T) {
	tt := []struct {
		description          string
		robotRequestBody     string
		expectedStatusCode   int
		expectedLocation     float64
		serviceErrorExpected bool
		expectedErrorMessage string
	}{
		{
			description: "databank located successfully",
			robotRequestBody: `{
				"x": "2.0",
				"y": "2.0",
				"z": "2.0",
				"vel": "5.0"
			}`,
			expectedStatusCode: http.StatusOK,
			expectedLocation:   11.0,
		},
		{
			description: "databank located successfully 2",
			robotRequestBody: `{
				"x": "123.12",
				"y": "456.56",
				"z": "789.89",
				"vel": "20.0"
			}`,
			expectedStatusCode: http.StatusOK,
			expectedLocation:   1389.57,
		},
		{
			description:        "empty robot status request body",
			robotRequestBody:   "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "invalid json request body",
			robotRequestBody:   "{invalid-json,",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:          "invalid robot status request body",
			robotRequestBody:     `{"a":"1"}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectedErrorMessage: databank.ErrInvalidRobotStatus.Error(),
		},
		{
			description: "databank location service failed",
			robotRequestBody: `{
				"x": "123.12",
				"y": "456.56",
				"z": "789.89",
				"vel": "20.0"
			}`,
			expectedStatusCode:   http.StatusInternalServerError,
			serviceErrorExpected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body := bytes.NewReader([]byte(tc.robotRequestBody))
			req, err := http.NewRequest(http.MethodPost, "", body)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			c.Request = req

			server, err := app.NewApp()
			defer server.Shutdown()
			if err != nil {
				t.Fatalf("could not create app: %v", err)
			}
			service := buildMockLocationService(tc.expectedLocation, tc.serviceErrorExpected)
			controller := databank.NewController(server.GetSectorID(), service)

			controller.LocatorHandler(c)

			res := w.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.expectedErrorMessage != "" {
				if res.StatusCode != tc.expectedStatusCode {
					t.Errorf("expected status code was: %d, got %d", tc.expectedStatusCode, res.StatusCode)
				}
				// check error
				var apiError databank.LocationAPIError
				err := json.Unmarshal(b, &apiError)
				if err != nil {
					t.Fatalf("invalid error returned from location handler, %s", b)
				}
				if apiError.Message != tc.expectedErrorMessage {
					t.Errorf("expected message %s, got %s", tc.expectedErrorMessage, apiError.Message)
				}
				return
			}
			if w.Code != tc.expectedStatusCode {
				t.Errorf("unexpected status code, want: %d, got: %d", tc.expectedStatusCode, w.Code)
			}
			var locationResponse databank.Location
			err = json.Unmarshal(b, &locationResponse)
			if err != nil {
				t.Fatalf("invalid location response returned from handler, %s", b)
			}
			if fmt.Sprintf("%f", locationResponse.Location) != fmt.Sprintf("%f", tc.expectedLocation) {
				t.Errorf("unexpected location returned, want: %f, got: %f", tc.expectedLocation, locationResponse.Location)
			}
		})
	}
}
