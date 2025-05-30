package http

import (
	"context"
	"encoding/json"
	"errors"
	"homework/internal/usecase"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	useCases    UseCases
	connections map[int64][]*websocket.Conn
	mu          sync.Mutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewWebSocketHandler(useCases UseCases) *WebSocketHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketHandler{
		useCases:    useCases,
		connections: make(map[int64][]*websocket.Conn),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context, id int64) error {
	_, err := h.useCases.Sensor.GetSensorByID(h.ctx, id)
	if err != nil {
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return err
		}
		return err
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}

	h.mu.Lock()
	h.connections[id] = append(h.connections[id], conn)
	h.mu.Unlock()

	ctx := conn.CloseRead(context.Background())
	go h.handleConnection(ctx, conn, id)

	return nil
}

func (h *WebSocketHandler) handleConnection(c context.Context, conn *websocket.Conn, sensorID int64) {
	defer func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		conn.Close(websocket.StatusNormalClosure, "connection closed")
		h.removeConnection(sensorID, conn)
	}()

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			event, err := h.useCases.Event.GetLastEventBySensorID(h.ctx, sensorID)
			if err != nil {
				if !errors.Is(err, usecase.ErrEventNotFound) {
					log.Printf("failed to get last event for sensor %d: %v", sensorID, err)
				}
				continue
			}

			msg, err := json.Marshal(event)
			if err != nil {
				log.Printf("failed to marshal event: %v", err)
				continue
			}

			err = conn.Write(h.ctx, websocket.MessageText, msg)
			if err != nil {
				log.Printf("failed to write to websocket: %v", err)
				return
			}

		case <-c.Done():
			return
		}
	}
}

func (h *WebSocketHandler) removeConnection(sensorID int64, conn *websocket.Conn) {
	conns := h.connections[sensorID]
	for i, c := range conns {
		if c == conn {
			h.connections[sensorID] = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}

func (h *WebSocketHandler) Shutdown() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.cancel()
	for _, conns := range h.connections {
		for _, conn := range conns {
			conn.Close(websocket.StatusGoingAway, "server shutdown")
		}
	}

	return nil
}
