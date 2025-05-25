package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// DeploymentMetrics represents metrics for a deployment
type DeploymentMetrics struct {
	ID           uuid.UUID           `json:"id" db:"id"`
	DeploymentID uuid.UUID       `json:"deployment_id" db:"deployment_id"`
	MetricDate   time.Time       `json:"metric_date" db:"metric_date"`
	TotalTrades  int             `json:"total_trades" db:"total_trades"`
	WinTrades    int             `json:"win_trades" db:"win_trades"`
	LossTrades   int             `json:"loss_trades" db:"loss_trades"`
	TotalPNL     float64         `json:"total_pnl" db:"total_pnl"`
	MaxDrawdown  float64         `json:"max_drawdown" db:"max_drawdown"`
	SharpeRatio  *float64        `json:"sharpe_ratio" db:"sharpe_ratio"`
	SortinoRatio *float64        `json:"sortino_ratio" db:"sortino_ratio"`
	WinRate      *float64        `json:"win_rate" db:"win_rate"`
	AvgWin       *float64        `json:"avg_win" db:"avg_win"`
	AvgLoss      *float64        `json:"avg_loss" db:"avg_loss"`
	Metadata     json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

// NewDeploymentMetrics creates new metrics for a deployment
func NewDeploymentMetrics(
	deploymentID uuid.UUID,
	metricDate time.Time,
	metadata map[string]interface{},
) (*DeploymentMetrics, error) {
	var metadataJSON json.RawMessage
	var err error

	if metadata != nil {
		metadataJSON, err = json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
	}

	return &DeploymentMetrics{
		DeploymentID: deploymentID,
		MetricDate:   metricDate,
		TotalTrades:  0,
		WinTrades:    0,
		LossTrades:   0,
		TotalPNL:     0,
		MaxDrawdown:  0,
		Metadata:     metadataJSON,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

// UpdateMetrics updates the metrics with trade results
func (m *DeploymentMetrics) UpdateMetrics(
	totalTrades int,
	winTrades int,
	lossTrades int,
	totalPNL float64,
	maxDrawdown float64,
) {
	m.TotalTrades = totalTrades
	m.WinTrades = winTrades
	m.LossTrades = lossTrades
	m.TotalPNL = totalPNL
	m.MaxDrawdown = maxDrawdown

	// Calculate derived metrics
	if totalTrades > 0 {
		winRate := float64(winTrades) / float64(totalTrades)
		m.WinRate = &winRate
	}

	if winTrades > 0 {
		avgWin := totalPNL / float64(winTrades)
		m.AvgWin = &avgWin
	}

	if lossTrades > 0 {
		avgLoss := totalPNL / float64(lossTrades)
		m.AvgLoss = &avgLoss
	}

	m.UpdatedAt = time.Now()
}

// SetRatios sets the Sharpe and Sortino ratios
func (m *DeploymentMetrics) SetRatios(sharpe, sortino float64) {
	m.SharpeRatio = &sharpe
	m.SortinoRatio = &sortino
	m.UpdatedAt = time.Now()
}
