package http

import (
	"encoding/json"
	"errors"
	"homework/internal/domain"
	"homework/internal/gateways/http/models"
	"homework/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

func sensorToDomain(sensorModel *models.SensorToCreate) (*domain.Sensor, error) {
	return &domain.Sensor{
		SerialNumber: *sensorModel.SerialNumber,
		Type:         domain.SensorType(*sensorModel.Type),
		Description:  *sensorModel.Description,
		IsActive:     *sensorModel.IsActive,
	}, nil
}

func sensorFromDomain(sensor *domain.Sensor) *models.Sensor {
	sensorType := string(sensor.Type)
	return &models.Sensor{
		ID:           &sensor.ID,
		SerialNumber: &sensor.SerialNumber,
		Type:         &sensorType,
		CurrentState: &sensor.CurrentState,
		Description:  &sensor.Description,
		IsActive:     &sensor.IsActive,
		RegisteredAt: (*strfmt.DateTime)(&sensor.RegisteredAt),
		LastActivity: (*strfmt.DateTime)(&sensor.LastActivity),
	}
}

func sensorsFromDomain(sensors []domain.Sensor) []*models.Sensor {
	var result []*models.Sensor
	for _, sensor := range sensors {
		result = append(result, sensorFromDomain(&sensor))
	}
	return result
}

func handleEventRegistration(c *gin.Context, eventUC *usecase.Event) {
	var eventModel models.SensorEvent
	if err := c.ShouldBindJSON(&eventModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := eventModel.Validate(strfmt.Default); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	event := domain.Event{
		SensorSerialNumber: *eventModel.SensorSerialNumber,
		Payload:            *eventModel.Payload,
		Timestamp:          time.Now(),
	}

	if err := eventUC.ReceiveEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func handleGetSensors(c *gin.Context, sensorUC *usecase.Sensor) {
	handleSensors(c, sensorUC, true)
}

func handleHeadSensors(c *gin.Context, sensorUC *usecase.Sensor) {
	handleSensors(c, sensorUC, false)
}

func handleSensors(c *gin.Context, sensorUC *usecase.Sensor, isGetReq bool) {
	sensors, err := sensorUC.GetSensors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseModels := sensorsFromDomain(sensors)
	setStatusOrJson(c, responseModels, isGetReq)
}

func handleSensorRegistration(c *gin.Context, sensorUC *usecase.Sensor) {
	var sensorModel models.SensorToCreate
	if err := c.ShouldBindJSON(&sensorModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := sensorModel.Validate(strfmt.Default); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	sensor, err := sensorToDomain(&sensorModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := sensorUC.RegisterSensor(c.Request.Context(), sensor)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	responseModel := sensorFromDomain(result)
	setStatusOrJson(c, responseModel, true)
}

func handleGetSensor(c *gin.Context, sensorUC *usecase.Sensor) {
	handleSensor(c, sensorUC, true)
}

func handleHeadSensor(c *gin.Context, sensorUC *usecase.Sensor) {
	handleSensor(c, sensorUC, false)
}

func handleSensor(c *gin.Context, sensorUC *usecase.Sensor, isGetReq bool) {
	sensorID, ok := parseAndValidateID(c, "sensor_id")
	if !ok {
		return
	}

	sensor, err := sensorUC.GetSensorByID(c.Request.Context(), sensorID)
	if err != nil {
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	responseModel := sensorFromDomain(sensor)
	setStatusOrJson(c, responseModel, isGetReq)
}

func handleGetUserSensors(c *gin.Context, userUC *usecase.User) {
	handleUserSensors(c, userUC, true)
}

func handleHeadUserSensors(c *gin.Context, userUC *usecase.User) {
	handleUserSensors(c, userUC, false)
}

func handleUserSensors(c *gin.Context, userUC *usecase.User, isGetReq bool) {
	userID, ok := parseAndValidateID(c, "user_id")
	if !ok {
		return
	}

	sensors, err := userUC.GetUserSensors(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	responseModels := sensorsFromDomain(sensors)
	setStatusOrJson(c, responseModels, isGetReq)
}

func handleUserRegistration(c *gin.Context, userUC *usecase.User) {
	var userModel models.UserToCreate
	if err := c.ShouldBindJSON(&userModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := userModel.Validate(strfmt.Default); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user := domain.User{
		Name: *userModel.Name,
	}

	result, err := userUC.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	responseModel := &models.User{
		ID:   &result.ID,
		Name: &result.Name,
	}

	setStatusOrJson(c, responseModel, true)
}

func handleBindSensorToUser(c *gin.Context, userUC *usecase.User) {
	userID, ok := parseAndValidateID(c, "user_id")
	if !ok {
		return
	}

	var bindingModel models.SensorToUserBinding
	if err := c.ShouldBindJSON(&bindingModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := bindingModel.Validate(strfmt.Default); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	binding := domain.SensorOwner{
		SensorID: *bindingModel.SensorID,
	}

	if err := userUC.AttachSensorToUser(c.Request.Context(), userID, binding.SensorID); err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) || errors.Is(err, usecase.ErrSensorNotFound) {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusCreated)
}

func handleGetLastEvent(c *gin.Context, wsh *WebSocketHandler) {
	sensorID, ok := parseAndValidateID(c, "sensor_id")
	if !ok {
		return
	}

	if err := wsh.Handle(c, sensorID); err != nil {
		if errors.Is(err, usecase.ErrSensorNotFound) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func setStatusOrJson(c *gin.Context, data interface{}, isJson bool) {
	if isJson {
		c.JSON(http.StatusOK, data)
	} else {
		jsonData, _ := json.Marshal(data)
		c.Header("Content-Length", strconv.Itoa(len(jsonData)))
		c.Status(http.StatusOK)
	}
}

func parseAndValidateID(c *gin.Context, param string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid " + param + " id"})
		return 0, false
	}
	return id, true
}
