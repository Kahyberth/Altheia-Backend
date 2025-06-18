package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid"
)

type Handler struct {
	hub     *Hub
	service *Service
}

func NewHandler(hub *Hub, service *Service) *Handler {
	return &Handler{
		hub:     hub,
		service: service,
	}
}

func (h *Handler) WebSocketUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		c.Locals("allowed", true)
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func (h *Handler) HandleWebSocket(c *websocket.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic en HandleWebSocket: %v", r)
		}
	}()

	clinicID := c.Query("clinic_id")
	if clinicID == "" {
		log.Println("clinic_id es requerido para la conexión WebSocket")
		c.Close()
		return
	}

	clientID, _ := gonanoid.Nanoid()

	client := &Client{
		ID:       clientID,
		Conn:     &FiberWebSocketConnection{conn: c},
		ClinicID: clinicID,
		Send:     make(chan []byte, 256),
	}

	log.Printf("Nueva conexión WebSocket establecida: %s para clínica: %s", clientID, clinicID)

	h.hub.Register <- client

	go h.handleWrite(client)
	h.handleRead(client)
}

func (h *Handler) handleWrite(client *Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		log.Printf("handleWrite terminado para cliente %s", client.ID)
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				log.Printf("Canal cerrado para cliente %s", client.ID)
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Error escribiendo mensaje al cliente %s: %v", client.ID, err)
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error enviando ping al cliente %s: %v", client.ID, err)
				return
			}
		}
	}
}

func (h *Handler) handleRead(client *Client) {
	defer func() {
		log.Printf("handleRead terminado para cliente %s", client.ID)
		h.hub.Unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		messageType, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error inesperado de WebSocket para cliente %s: %v", client.ID, err)
			} else {
				log.Printf("Cliente %s se desconectó normalmente", client.ID)
			}
			break
		}

		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		if messageType == websocket.TextMessage {
			log.Printf("Mensaje recibido del cliente %s: %s", client.ID, string(message))
		}
	}
}

func (h *Handler) GetStats(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"connected_clients": h.hub.GetConnectedClients(),
		"timestamp":         time.Now(),
	})
}

func (h *Handler) GetClinicStats(c *fiber.Ctx) error {
	clinicID := c.Params("clinicId")
	if clinicID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "clinic_id es requerido"})
	}

	return c.JSON(fiber.Map{
		"clinic_id":         clinicID,
		"connected_clients": h.hub.GetClinicClients(clinicID),
		"timestamp":         time.Now(),
	})
}

func (h *Handler) BroadcastMessage(c *fiber.Ctx) error {
	var request struct {
		Type     string      `json:"type"`
		Data     interface{} `json:"data"`
		ClinicID string      `json:"clinic_id,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	if request.ClinicID != "" {
		h.service.BroadcastToClinic(request.ClinicID, request.Type, request.Data)
	} else {
		message := Message{
			Type:      request.Type,
			Data:      request.Data,
			Timestamp: time.Now(),
		}

		if data, err := json.Marshal(message); err == nil {
			h.hub.Broadcast <- data
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Error marshaling message"})
		}
	}

	return c.JSON(fiber.Map{"message": "Mensaje enviado exitosamente"})
}

type FiberWebSocketConnection struct {
	conn *websocket.Conn
}

func (f *FiberWebSocketConnection) WriteMessage(messageType int, data []byte) error {
	return f.conn.WriteMessage(messageType, data)
}

func (f *FiberWebSocketConnection) ReadMessage() (messageType int, p []byte, err error) {
	return f.conn.ReadMessage()
}

func (f *FiberWebSocketConnection) Close() error {
	return f.conn.Close()
}

func (f *FiberWebSocketConnection) SetReadDeadline(t time.Time) error {
	return f.conn.SetReadDeadline(t)
}

func (f *FiberWebSocketConnection) SetPongHandler(h func(appData string) error) {
	f.conn.SetPongHandler(h)
}

func (f *FiberWebSocketConnection) SetWriteDeadline(t time.Time) error {
	return f.conn.SetWriteDeadline(t)
}
