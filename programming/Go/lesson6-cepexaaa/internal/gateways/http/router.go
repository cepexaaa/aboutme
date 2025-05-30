package http

import (
	"homework/internal/usecase"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type anyCases interface {
	*usecase.Event | *usecase.User | *usecase.Sensor | *WebSocketHandler
}

func setupRouter(r *gin.Engine, uc UseCases, wsh *WebSocketHandler) {
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "route not found"})
	})

	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	})

	eventsGroup := r.Group("/events")
	{
		registerHandler(eventsGroup, http.MethodPost, "", handleEventRegistration, uc.Event)
		eventsGroup.OPTIONS("", handleOptions([]string{"POST", "OPTIONS"}))
	}

	sensorsGroup := r.Group("/sensors")
	{
		registerHandler(sensorsGroup, http.MethodGet, "", handleGetSensors, uc.Sensor)
		registerHandler(sensorsGroup, http.MethodHead, "", handleHeadSensors, uc.Sensor)
		registerHandler(sensorsGroup, http.MethodPost, "", handleSensorRegistration, uc.Sensor)
		sensorsGroup.OPTIONS("", handleOptions([]string{"GET", "HEAD", "POST", "OPTIONS"}))

		registerHandler(sensorsGroup, http.MethodGet, "/:sensor_id", handleGetSensor, uc.Sensor)
		registerHandler(sensorsGroup, http.MethodHead, "/:sensor_id", handleHeadSensor, uc.Sensor)
		sensorsGroup.OPTIONS("/:sensor_id", handleOptions([]string{"GET", "HEAD", "OPTIONS"}))
		registerHandler(sensorsGroup, http.MethodGet, "/:sensor_id/events", handleGetLastEvent, wsh)
	}

	usersGroup := r.Group("/users")
	{
		registerHandler(usersGroup, http.MethodPost, "", handleUserRegistration, uc.User)
		usersGroup.OPTIONS("", handleOptions([]string{"POST", "OPTIONS"}))

		userSensorsGroup := usersGroup.Group("/:user_id/sensors")
		{
			registerHandler(userSensorsGroup, http.MethodGet, "", handleGetUserSensors, uc.User)
			registerHandler(userSensorsGroup, http.MethodHead, "", handleHeadUserSensors, uc.User)
			registerHandler(userSensorsGroup, http.MethodPost, "", handleBindSensorToUser, uc.User)
			userSensorsGroup.OPTIONS("", handleOptions([]string{"GET", "HEAD", "POST", "OPTIONS"}))
		}
	}
}

func handleOptions(methods []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	}
}

func registerHandler[T anyCases](group *gin.RouterGroup, method, path string, handler func(*gin.Context, T), uc T) {
	methodHandlers := map[string]struct {
		register  func(string, ...gin.HandlerFunc) gin.IRoutes
		validator gin.HandlerFunc
	}{
		http.MethodGet:  {group.GET, validateAccept},
		http.MethodPost: {group.POST, validateContentType},
		http.MethodHead: {group.HEAD, validateAccept},
	}

	registry, ok := methodHandlers[method]
	if !ok {
		panic("unsupported http method: " + method)
	}

	isWebSocket := reflect.TypeOf(uc) == reflect.TypeOf((*WebSocketHandler)(nil))
	if isWebSocket && method == http.MethodGet {
		registry.register(path, wrapHandler(handler, uc))
	} else {
		registry.register(path, registry.validator, wrapHandler(handler, uc))
	}
}

func wrapHandler[T anyCases](handler func(*gin.Context, T), uc T) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, uc)
	}
}

func validateContentType(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "unsupported media type"})
		c.Abort()
	}
}

func validateAccept(c *gin.Context) {
	if c.GetHeader("Accept") != "application/json" {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "not acceptable"})
		c.Abort()
	}
}
