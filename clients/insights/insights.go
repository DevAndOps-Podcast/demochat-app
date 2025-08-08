package insights

import (
	"context"
	"demochat/config"
	"log"
	"net/http"
	"time"

	"github.com/imroc/req/v3"
)

type Client struct {
	*req.Client
}

type Insights struct {
	MostActiveUserID   int64   `json:"most_active_user_id"`
	TotalMessages      int64   `json:"total_messages"`
	AverageMessageRate float64 `json:"average_message_rate"`
}

func New(cfg *config.Config) Client {
	c := req.NewClient()

	c.BaseURL = cfg.InsightsService.BaseUrl
	c.SetTimeout(60 * time.Second)
	c.Headers = http.Header{}
	c.Headers.Add("Content-Type", "application/json")
	c.Headers.Add("service-secret", cfg.InsightsService.ApiKey)

	return Client{c}
}

type PublishMessageRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
}

func (c Client) PublishMessage(ctx context.Context, data PublishMessageRequest) error {
	resp, err := c.R().
		SetBody(data).
		SetRetryCount(1).
		Post("/messages")

	if resp.IsErrorState() {
		log.Println("call to /messages failed.", resp.Status)
		return resp.Err
	}

	return err
}

func (c Client) ChatRoomInsights(ctx context.Context) Insights {
	var insights Insights
	c.R().
		SetSuccessResult(&insights).
		Get("/messages")

	return insights
}
