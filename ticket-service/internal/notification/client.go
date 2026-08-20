package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}
type NotificationServer struct {
	BaseURL    string
	HTTPClient *http.Client
}
type NotificationRequest struct {
	EventID         string    `json:"event_id"`
	TicketID        int       `json:"ticket_id"`
	ActorUserID     int       `json:"actor_user_id"`
	RecipientUserID int       `json:"recipient_user_id"`
	Type            string    `json:"type"`
	Message         string    `json:"message"`
	OccurredAt      time.Time `json:"occurred_at"`
}
type NotificationError struct {
	StatusCode int
}

func NewClient(baseURL string) *Client { //burada istegimiz url için client oluşturyoruz eğer sonra başka servis eklenirse kolay implement edilebilir
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}
func (c *Client) CreateNotification(notification NotificationRequest) error {
	data, err := json.Marshal(notification) //buradaki Marshal komutu struct olan requsti şablonu kllanrak jsona çeviriyor

	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.BaseURL+"/notifications", bytes.NewBuffer(data)) //önce requsestimzi hazırladık detaylı bir şekilde biligleri ekliyoruz  authantication eklenebilir ileride
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req) //isteği burada gönderdik
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("error closing response body: %s", err)
		}
	}() //bağlantıyı ,isteği kapatıyoruz
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &NotificationError{
			StatusCode: resp.StatusCode,
		}
	}
	return nil
}
func (e *NotificationError) Error() string {
	return fmt.Sprintf(
		"notification service returned status: %d",
		e.StatusCode,
	)
}
