package dto

import (
	"encoding/json"
	"time"

	"algo/internal/modules/deployment/domain/model"

	"github.com/google/uuid"
)

// DeploymentPositionLogDTO represents a data transfer object for deployment position logs
type DeploymentPositionLogDTO struct {
	ID                   uuid.UUID       `json:"id"`
	DeploymentPositionID int64           `json:"deployment_position_id"`
	Results              json.RawMessage `json:"results"`
	Timestamp            time.Time       `json:"timestamp"`
	CreatedAt            time.Time       `json:"created_at"`
}

// ToDTO converts a model to a DTO
func (dto *DeploymentPositionLogDTO) ToDTO(model *model.DeploymentPositionLog) *DeploymentPositionLogDTO {
	return &DeploymentPositionLogDTO{
		ID:                   model.ID,
		DeploymentPositionID: model.DeploymentPositionID,
		Results:              model.Results,
		Timestamp:            model.Timestamp,
		CreatedAt:            model.CreatedAt,
	}
}

// ToModel converts a DTO to a model
func (dto *DeploymentPositionLogDTO) ToModel() *model.DeploymentPositionLog {
	return &model.DeploymentPositionLog{
		ID:                   dto.ID,
		DeploymentPositionID: dto.DeploymentPositionID,
		Results:              dto.Results,
		Timestamp:            dto.Timestamp,
		CreatedAt:            dto.CreatedAt,
	}
}

// ToDTOList converts a list of models to a list of DTOs
func ToDTOList(models []*model.DeploymentPositionLog) []*DeploymentPositionLogDTO {
	dtos := make([]*DeploymentPositionLogDTO, len(models))
	for i, model := range models {
		dto := &DeploymentPositionLogDTO{}
		dtos[i] = dto.ToDTO(model)
	}
	return dtos
}
