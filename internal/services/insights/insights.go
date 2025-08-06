package insights

import (
	"context"
	"demochat/clients/insights"

	"go.uber.org/fx"
)

type Service interface {
	GetInsights(ctx context.Context) insights.Insights
}

type service struct {
	insighsClient insights.Client
}

type Params struct {
	fx.In

	InsightsClient insights.Client
}

func New(p Params) Service {
	return &service{
		insighsClient: p.InsightsClient,
	}
}

func (s *service) GetInsights(ctx context.Context) insights.Insights {
	return s.insighsClient.ChatRoomInsights(ctx)
}
