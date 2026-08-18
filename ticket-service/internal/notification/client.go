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
	TicketID int    `json:"ticket_id"`
	UserID   int    `json:"user_id"`
	Type     string `json:"type"`
	Message  string `json:"message"`
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { //200 ile 300 arasındaki kodlar genelde onaylandı anlamında taşır zaten biz birkaçını kullnadık .diğer kodlarda 500,404 gibi yi kapsayacak her biri için case yazmaya gerek yok
		return fmt.Errorf(
			"notification service returned status: %d",
			resp.StatusCode,
		)
	}
	return nil
}
