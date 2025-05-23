package dto

import (
	"encoding/json"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// CreateMetricsRequest is the request for creating deployment metrics
type CreateMetricsRequest struct {
	DeploymentID uuid.UUID              `json:"deployment_id" validate:"required"`
	MetricDate   time.Time              `json:"metric_date" validate:"required"`
	TotalTrades  int                    `json:"total_trades" validate:"gte=0"`
	WinTrades    int                    `json:"win_trades" validate:"gte=0"`
	LossTrades   int                    `json:"loss_trades" validate:"gte=0"`
	TotalPNL     float64                `json:"total_pnl"`
	MaxDrawdown  float64                `json:"max_drawdown" validate:"gte=0"`
	SharpeRatio  *float64               `json:"sharpe_ratio,omitempty"`
	SortinoRatio *float64               `json:"sortino_ratio,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateMetricsRequest is the request for updating deployment metrics
type UpdateMetricsRequest struct {
	TotalTrades  int                    `json:"total_trades" validate:"gte=0"`
	WinTrades    int                    `json:"win_trades" validate:"gte=0"`
	LossTrades   int                    `json:"loss_trades" validate:"gte=0"`
	TotalPNL     float64                `json:"total_pnl"`
	MaxDrawdown  float64                `json:"max_drawdown" validate:"gte=0"`
	SharpeRatio  *float64               `json:"sharpe_ratio,omitempty"`
	SortinoRatio *float64               `json:"sortino_ratio,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// MetricsFilter is used to filter metrics queries
type MetricsFilter struct {
	DeploymentID uuid.UUID  `json:"deployment_id" validate:"required"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
}

// MetricsResponse is the response for deployment metrics
type MetricsResponse struct {
	ID           int64           `json:"id"`
	DeploymentID uuid.UUID       `json:"deployment_id"`
	MetricDate   time.Time       `json:"metric_date"`
	TotalTrades  int             `json:"total_trades"`
	WinTrades    int             `json:"win_trades"`
	LossTrades   int             `json:"loss_trades"`
	TotalPNL     float64         `json:"total_pnl"`
	MaxDrawdown  float64         `json:"max_drawdown"`
	SharpeRatio  *float64        `json:"sharpe_ratio,omitempty"`
	SortinoRatio *float64        `json:"sortino_ratio,omitempty"`
	WinRate      *float64        `json:"win_rate,omitempty"`
	AvgWin       *float64        `json:"avg_win,omitempty"`
	AvgLoss      *float64        `json:"avg_loss,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ToMetricsResponse converts a model.DeploymentMetrics to dto.MetricsResponse
func ToMetricsResponse(metrics *model.DeploymentMetrics) *MetricsResponse {
	return &MetricsResponse{
		ID:           metrics.ID,
		DeploymentID: metrics.DeploymentID,
		MetricDate:   metrics.MetricDate,
		TotalTrades:  metrics.TotalTrades,
		WinTrades:    metrics.WinTrades,
		LossTrades:   metrics.LossTrades,
		TotalPNL:     metrics.TotalPNL,
		MaxDrawdown:  metrics.MaxDrawdown,
		SharpeRatio:  metrics.SharpeRatio,
		SortinoRatio: metrics.SortinoRatio,
		WinRate:      metrics.WinRate,
		AvgWin:       metrics.AvgWin,
		AvgLoss:      metrics.AvgLoss,
		Metadata:     metrics.Metadata,
		CreatedAt:    metrics.CreatedAt,
		UpdatedAt:    metrics.UpdatedAt,
	}
}

// ToMetricsResponses converts a slice of model.DeploymentMetrics to dto.MetricsResponse
func ToMetricsResponses(metrics []*model.DeploymentMetrics) []*MetricsResponse {
	responses := make([]*MetricsResponse, len(metrics))
	for i, m := range metrics {
		responses[i] = ToMetricsResponse(m)
	}
	return responses
}
