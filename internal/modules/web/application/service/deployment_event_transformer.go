package service

import (
	"fmt"
	"time"

	"algo/internal/modules/deployment/application/dto"
)

// DeploymentEventTransformer transforms deployment events into page-specific web events
type DeploymentEventTransformer struct{}

// NewDeploymentEventTransformer creates a new deployment event transformer
func NewDeploymentEventTransformer() *DeploymentEventTransformer {
	return &DeploymentEventTransformer{}
}

// TransformToWebEvents transforms a deployment event into one or more page-specific web events
func (t *DeploymentEventTransformer) TransformToWebEvents(deploymentEvent *dto.DeploymentEvent) []*dto.WebEvent {
	var webEvents []*dto.WebEvent

	// Common metadata for all events
	baseMetadata := map[string]interface{}{
		"user_id":       deploymentEvent.Metadata.UserID,
		"deployment_id": deploymentEvent.DeploymentID.String(),
		"strategy_name": deploymentEvent.Metadata.StrategyName,
	}

	switch deploymentEvent.EventType {
	case dto.EventTypePositionOpened:
		webEvents = append(webEvents, t.transformPositionOpened(deploymentEvent, baseMetadata)...)
	case dto.EventTypePositionClosed:
		webEvents = append(webEvents, t.transformPositionClosed(deploymentEvent, baseMetadata)...)
	case dto.EventTypeEntrySignalDetected:
		webEvents = append(webEvents, t.transformSignalDetected(deploymentEvent, baseMetadata, "entry")...)
	case dto.EventTypeExitSignalDetected:
		webEvents = append(webEvents, t.transformSignalDetected(deploymentEvent, baseMetadata, "exit")...)
	case dto.EventTypeDeploymentActivated:
		webEvents = append(webEvents, t.transformDeploymentStateChanged(deploymentEvent, baseMetadata, "activated")...)
	case dto.EventTypeDeploymentDeactivated:
		webEvents = append(webEvents, t.transformDeploymentStateChanged(deploymentEvent, baseMetadata, "deactivated")...)
	case dto.EventTypeDeploymentError:
		webEvents = append(webEvents, t.transformDeploymentError(deploymentEvent, baseMetadata)...)
	case dto.EventTypeStopLossTriggered:
		webEvents = append(webEvents, t.transformStopLossTriggered(deploymentEvent, baseMetadata)...)
	case dto.EventTypeTargetProfitTriggered:
		webEvents = append(webEvents, t.transformTargetProfitTriggered(deploymentEvent, baseMetadata)...)
	}

	return webEvents
}

// transformPositionOpened transforms position opened events
func (t *DeploymentEventTransformer) transformPositionOpened(event *dto.DeploymentEvent, baseMetadata map[string]interface{}) []*dto.WebEvent {
	var events []*dto.WebEvent

	// Extract position data (assuming it's in the correct format)
	positionData, ok := event.Data.(*dto.PositionOpenedData)
	if !ok {
		return events
	}

	// Dashboard event - recent activity
	dashboardEvent := &dto.WebEvent{
		EventType: "dashboard.recent_activity",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityMedium,
		Data: dto.DashboardRecentActivityData{
			Activity:       fmt.Sprintf("Position opened: %s %.0f @ %.2f", positionData.Direction, positionData.Quantity, positionData.EntryPrice),
			DeploymentID:   event.DeploymentID.String(),
			DeploymentName: event.Metadata.StrategyName,
			Timestamp:      event.Timestamp,
			Type:           "position_opened",
			Details: map[string]interface{}{
				"direction":   positionData.Direction,
				"quantity":    positionData.Quantity,
				"entry_price": positionData.EntryPrice,
				"position_id": positionData.PositionID,
			},
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, dashboardEvent)

	// Deployment list event
	deploymentListEvent := &dto.WebEvent{
		EventType: "deployment_list.item_updated",
		PageType:  dto.PageTypeDeploymentList,
		Priority:  dto.EventPriorityMedium,
		Data: dto.DeploymentListItemData{
			DeploymentID:  event.DeploymentID.String(),
			Name:          event.Metadata.StrategyName,
			Status:        "active",
			LastSignal:    "position_opened",
			LastActive:    event.Timestamp,
			OpenPositions: 1, // This would need to be calculated
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, deploymentListEvent)

	// Deployment detail event
	deploymentDetailEvent := &dto.WebEvent{
		EventType: "deployment_detail.position_opened",
		PageType:  dto.PageTypeDeploymentDetail,
		Priority:  dto.EventPriorityHigh,
		Data: dto.DeploymentDetailPositionData{
			PositionID: positionData.PositionID,
			EntryPrice: positionData.EntryPrice,
			Quantity:   int(positionData.Quantity),
			Direction:  positionData.Direction,
			Status:     "open",
			Timestamp:  event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, deploymentDetailEvent)

	return events
}

// transformPositionClosed transforms position closed events
func (t *DeploymentEventTransformer) transformPositionClosed(event *dto.DeploymentEvent, baseMetadata map[string]interface{}) []*dto.WebEvent {
	var events []*dto.WebEvent

	positionData, ok := event.Data.(*dto.PositionClosedData)
	if !ok {
		return events
	}

	pnlColor := "success"
	if positionData.PNL < 0 {
		pnlColor = "error"
	}

	// Dashboard event
	dashboardEvent := &dto.WebEvent{
		EventType: "dashboard.recent_activity",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityMedium,
		Data: dto.DashboardRecentActivityData{
			Activity:       fmt.Sprintf("Position closed: P&L %.2f (%s)", positionData.PNL, positionData.ExitReason),
			DeploymentID:   event.DeploymentID.String(),
			DeploymentName: event.Metadata.StrategyName,
			Timestamp:      event.Timestamp,
			Type:           "position_closed",
			Details: map[string]interface{}{
				"pnl":         positionData.PNL,
				"exit_reason": positionData.ExitReason,
				"exit_price":  positionData.ExitPrice,
				"position_id": positionData.PositionID,
				"color":       pnlColor,
			},
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, dashboardEvent)

	// Deployment detail event
	deploymentDetailEvent := &dto.WebEvent{
		EventType: "deployment_detail.position_closed",
		PageType:  dto.PageTypeDeploymentDetail,
		Priority:  dto.EventPriorityHigh,
		Data: dto.DeploymentDetailPositionData{
			PositionID: positionData.PositionID,
			ExitPrice:  &positionData.ExitPrice,
			PnL:        &positionData.PNL,
			Status:     "closed",
			Timestamp:  event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, deploymentDetailEvent)

	// Portfolio balance update
	portfolioEvent := &dto.WebEvent{
		EventType: "portfolio.balance_updated",
		PageType:  dto.PageTypePortfolio,
		Priority:  dto.EventPriorityMedium,
		Data: dto.PortfolioBalanceData{
			RealizedPnL: positionData.PNL,
			LastUpdated: event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, portfolioEvent)

	return events
}

// transformSignalDetected transforms signal detection events
func (t *DeploymentEventTransformer) transformSignalDetected(event *dto.DeploymentEvent, baseMetadata map[string]interface{}, signalType string) []*dto.WebEvent {
	var events []*dto.WebEvent

	// Use SignalDetectedData for both entry and exit signals
	signalData, ok := event.Data.(*dto.SignalDetectedData)
	if !ok {
		return events
	}

	// Deployment detail event for signals
	deploymentDetailEvent := &dto.WebEvent{
		EventType: "deployment_detail.signal_detected",
		PageType:  dto.PageTypeDeploymentDetail,
		Priority:  dto.EventPriorityHigh,
		Data: dto.DeploymentDetailSignalData{
			SignalType: signalType,
			Price:      signalData.MarketPrice,
			Confidence: 1.0, // Default confidence since it's not in the original data
			Reason:     signalData.Expression,
			Timestamp:  event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, deploymentDetailEvent)

	// Dashboard activity
	dashboardEvent := &dto.WebEvent{
		EventType: "dashboard.recent_activity",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityLow,
		Data: dto.DashboardRecentActivityData{
			Activity:       fmt.Sprintf("%s signal detected", signalType),
			DeploymentID:   event.DeploymentID.String(),
			DeploymentName: event.Metadata.StrategyName,
			Timestamp:      event.Timestamp,
			Type:           fmt.Sprintf("%s_signal_detected", signalType),
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, dashboardEvent)

	return events
}

// transformDeploymentStateChanged transforms deployment state change events
func (t *DeploymentEventTransformer) transformDeploymentStateChanged(event *dto.DeploymentEvent, baseMetadata map[string]interface{}, newState string) []*dto.WebEvent {
	var events []*dto.WebEvent

	// Dashboard summary update
	dashboardEvent := &dto.WebEvent{
		EventType: "dashboard.summary_updated",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityMedium,
		Data: dto.DashboardSummaryData{
			LastUpdated: event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, dashboardEvent)

	// Deployment list update
	deploymentListEvent := &dto.WebEvent{
		EventType: "deployment_list.item_updated",
		PageType:  dto.PageTypeDeploymentList,
		Priority:  dto.EventPriorityMedium,
		Data: dto.DeploymentListItemData{
			DeploymentID: event.DeploymentID.String(),
			Name:         event.Metadata.StrategyName,
			Status:       newState,
			LastActive:   event.Timestamp,
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, deploymentListEvent)

	// Recent activity
	activityEvent := &dto.WebEvent{
		EventType: "dashboard.recent_activity",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityLow,
		Data: dto.DashboardRecentActivityData{
			Activity:       fmt.Sprintf("Deployment %s", newState),
			DeploymentID:   event.DeploymentID.String(),
			DeploymentName: event.Metadata.StrategyName,
			Timestamp:      event.Timestamp,
			Type:           fmt.Sprintf("deployment_%s", newState),
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, activityEvent)

	return events
}

// transformDeploymentError transforms deployment error events
func (t *DeploymentEventTransformer) transformDeploymentError(event *dto.DeploymentEvent, baseMetadata map[string]interface{}) []*dto.WebEvent {
	var events []*dto.WebEvent

	errorData, ok := event.Data.(*dto.ErrorEventData)
	if !ok {
		return events
	}

	// High priority alert for all pages
	alertEvent := &dto.WebEvent{
		EventType: "alerts.new_notification",
		PageType:  "", // Empty means all pages
		Priority:  dto.EventPriorityHigh,
		Data: dto.AlertData{
			AlertID:        fmt.Sprintf("error_%d", time.Now().Unix()),
			Type:           "deployment_error",
			Severity:       "critical",
			Title:          "Deployment Error",
			Message:        fmt.Sprintf("Error in %s: %s", event.Metadata.StrategyName, errorData.ErrorMessage),
			RequiresAction: true,
			Timestamp:      event.Timestamp,
			Context: map[string]interface{}{
				"deployment_id": event.DeploymentID.String(),
				"error_code":    errorData.ErrorCode,
			},
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
		ExpiresAt: &[]time.Time{event.Timestamp.Add(24 * time.Hour)}[0],
	}
	events = append(events, alertEvent)

	// Dashboard alerts count update
	dashboardAlertsEvent := &dto.WebEvent{
		EventType: "dashboard.alerts_updated",
		PageType:  dto.PageTypeDashboard,
		Priority:  dto.EventPriorityHigh,
		Data: dto.DashboardAlertsData{
			CriticalCount: 1, // This would need to be calculated from actual alert counts
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, dashboardAlertsEvent)

	return events
}

// transformStopLossTriggered transforms stop loss triggered events
func (t *DeploymentEventTransformer) transformStopLossTriggered(event *dto.DeploymentEvent, baseMetadata map[string]interface{}) []*dto.WebEvent {
	var events []*dto.WebEvent

	stopLossData, ok := event.Data.(*dto.RiskEventData)
	if !ok {
		return events
	}

	// High priority notification
	alertEvent := &dto.WebEvent{
		EventType: "alerts.new_notification",
		PageType:  "",
		Priority:  dto.EventPriorityHigh,
		Data: dto.AlertData{
			AlertID:   fmt.Sprintf("stoploss_%d", time.Now().Unix()),
			Type:      "stop_loss",
			Severity:  "warning",
			Title:     "Stop Loss Triggered",
			Message:   fmt.Sprintf("Stop loss hit at %.2f for %s", stopLossData.TriggerPrice, event.Metadata.StrategyName),
			Timestamp: event.Timestamp,
			Context: map[string]interface{}{
				"position_id":   stopLossData.PositionID,
				"trigger_price": stopLossData.TriggerPrice,
				"deployment_id": event.DeploymentID.String(),
			},
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, alertEvent)

	return events
}

// transformTargetProfitTriggered transforms target profit triggered events
func (t *DeploymentEventTransformer) transformTargetProfitTriggered(event *dto.DeploymentEvent, baseMetadata map[string]interface{}) []*dto.WebEvent {
	var events []*dto.WebEvent

	targetData, ok := event.Data.(*dto.RiskEventData)
	if !ok {
		return events
	}

	// Success notification
	alertEvent := &dto.WebEvent{
		EventType: "alerts.new_notification",
		PageType:  "",
		Priority:  dto.EventPriorityMedium,
		Data: dto.AlertData{
			AlertID:   fmt.Sprintf("target_%d", time.Now().Unix()),
			Type:      "target_profit",
			Severity:  "info",
			Title:     "Target Reached",
			Message:   fmt.Sprintf("Target profit hit at %.2f for %s", targetData.TriggerPrice, event.Metadata.StrategyName),
			Timestamp: event.Timestamp,
			Context: map[string]interface{}{
				"position_id":   targetData.PositionID,
				"trigger_price": targetData.TriggerPrice,
				"deployment_id": event.DeploymentID.String(),
			},
		},
		Metadata:  baseMetadata,
		Timestamp: event.Timestamp,
	}
	events = append(events, alertEvent)

	return events
}
