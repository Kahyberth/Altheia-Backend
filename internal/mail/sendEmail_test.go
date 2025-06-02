package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment variables
	os.Setenv("EMAIL_USERNAME", "test@example.com")
	os.Setenv("EMAIL_PASSWORD", "testpassword")
	os.Setenv("EMAIL_HOST", "smtp.example.com")

	code := m.Run()
	os.Exit(code)
}

// Mock SMTP client for testing
type mockSMTPClient struct {
	from    string
	to      []string
	msg     []byte
	auth    smtp.Auth
	addr    string
	success bool
}

func (m *mockSMTPClient) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	m.from = from
	m.to = to
	m.msg = msg
	m.auth = auth
	m.addr = addr
	if !m.success {
		return fmt.Errorf("mock SMTP error")
	}
	return nil
}

func TestSendWelcomeMessage(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		to        []string
		mockSMTP  *mockSMTPClient
		wantErr   bool
		checkFunc func(*mockSMTPClient) error
	}{
		{
			name:     "Valid Welcome Message",
			username: "John Doe",
			to:       []string{"john@example.com"},
			mockSMTP: &mockSMTPClient{success: true},
			wantErr:  false,
			checkFunc: func(m *mockSMTPClient) error {
				if m.from == "" {
					return fmt.Errorf("expected from address, got empty")
				}
				if len(m.to) != 1 || m.to[0] != "john@example.com" {
					return fmt.Errorf("expected to address 'john@example.com', got %v", m.to)
				}
				if len(m.msg) == 0 {
					return fmt.Errorf("expected message content, got empty")
				}
				return nil
			},
		},
		{
			name:     "Empty Username",
			username: "",
			to:       []string{"john@example.com"},
			mockSMTP: &mockSMTPClient{success: true},
			wantErr:  true,
		},
		{
			name:     "Empty To Address",
			username: "John Doe",
			to:       []string{""},
			mockSMTP: &mockSMTPClient{success: true},
			wantErr:  true,
		},
		{
			name:     "SMTP Error",
			username: "John Doe",
			to:       []string{"john@example.com"},
			mockSMTP: &mockSMTPClient{success: false},
			wantErr:  false, // SMTP errors are logged but don't return error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original SMTP client and restore after test
			originalClient := smtpClient
			defer func() { smtpClient = originalClient }()

			// Set mock SMTP client
			SetSMTPClient(tt.mockSMTP)

			err := SendWelcomeMessage(tt.username, tt.to)

			if (err != nil) != tt.wantErr {
				t.Errorf("SendWelcomeMessage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkFunc != nil {
				if err := tt.checkFunc(tt.mockSMTP); err != nil {
					t.Errorf("SendWelcomeMessage() validation failed: %v", err)
				}
			}
		})
	}
}

func TestEmailTemplate(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     string
	}{
		{
			name:     "Valid Username",
			username: "John Doe",
			want:     "John Doe", // Template should contain the username
		},
		{
			name:     "Empty Username",
			username: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmailTemplate(tt.username)
			if tt.username != "" && !strings.Contains(got, tt.want) {
				t.Errorf("EmailTemplate() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}
