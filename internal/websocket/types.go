package websocket

import "time"

type ChartData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type PatientStatsMessage struct {
	Type string `json:"type"`
	Data struct {
		AgeDistribution    []ChartData `json:"age_distribution"`
		GenderDistribution []ChartData `json:"gender_distribution"`
		TotalPatients      int         `json:"total_patients"`
		Timestamp          time.Time   `json:"timestamp"`
	} `json:"data"`
}

type AppointmentStatsMessage struct {
	Type string `json:"type"`
	Data struct {
		AppointmentsByStatus []ChartData `json:"appointments_by_status"`
		AppointmentsByMonth  []ChartData `json:"appointments_by_month"`
		TotalAppointments    int         `json:"total_appointments"`
		TodayAppointments    int         `json:"today_appointments"`
		Timestamp            time.Time   `json:"timestamp"`
	} `json:"data"`
}

type ConsultationStatsMessage struct {
	Type string `json:"type"`
	Data struct {
		ConsultationsCreated []ChartData `json:"consultations_created"`
		TotalConsultations   int         `json:"total_consultations"`
		WeeklyConsultations  int         `json:"weekly_consultations"`
		Timestamp            time.Time   `json:"timestamp"`
	} `json:"data"`
}

type Client struct {
	ID       string
	Conn     WebSocketConnection
	ClinicID string
	Send     chan []byte
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	ClinicID  string      `json:"clinic_id"`
}

type WebSocketConnection interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	Close() error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
}
