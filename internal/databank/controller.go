package databank

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Controller abstracts the databank's location controller.
type Controller interface {
	LocatorHandler(c *gin.Context)
}

// NewController builds a databank's location controller.
func NewController(sectorID int, locationService Service) Controller {
	return &LocatorController{
		sectorID:        sectorID,
		locationService: locationService,
	}
}

// LocatorController holds the location services.
type LocatorController struct {
	sectorID        int
	locationService Service
}

// LocatorHandler receibes requests from robots and calculate the nearest databank.
func (l *LocatorController) LocatorHandler(c *gin.Context) {
	var robotState RobotState

	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()

	if err != nil {
		log.Errorf("[event: invalid_request_body][service: locator_controller] %s", err.Error())
		c.JSON(http.StatusBadRequest, newLocationAPIError(http.StatusBadRequest, err.Error()))
		return
	}
	if err := json.Unmarshal(body, &robotState); err != nil {
		log.Errorf("[event: invalid_request_body][service: locator_controller] %s", err.Error())
		c.JSON(http.StatusBadRequest, newLocationAPIError(http.StatusBadRequest, err.Error()))
		return
	}

	databankLocation, err := l.locationService.GetNearestDatabankLocation(robotState, l.sectorID)
	if err != nil {
		log.Errorf("[event: location_service_failed][service: locator_controller] %s", err.Error())
		c.JSON(http.StatusInternalServerError, newLocationAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, databankLocation)
}
