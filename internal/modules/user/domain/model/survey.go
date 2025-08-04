package model

import (
	"time"

	"github.com/google/uuid"
)

// NatureOfWork represents the nature of work enum
type NatureOfWork string

const (
	NatureOfWorkIndividual             NatureOfWork = "Individual"
	NatureOfWorkCompany                NatureOfWork = "Company"
	NatureOfWorkProprietaryTradingFirm NatureOfWork = "Proprietary_Trading_Firm"
	NatureOfWorkAssetManagement        NatureOfWork = "Asset_Management"
	NatureOfWorkFamilyOffice           NatureOfWork = "Family_Office"
	NatureOfWorkOther                  NatureOfWork = "Other"
)

// TradingType represents the trading type enum
type TradingType string

const (
	TradingTypeEquity TradingType = "Equity"
	TradingTypeFNO    TradingType = "FNO"
	TradingTypeBoth   TradingType = "Both"
)

// ExperienceLevel represents the experience level enum
type ExperienceLevel string

const (
	ExperienceLevelBeginner     ExperienceLevel = "Beginner"
	ExperienceLevelIntermediate ExperienceLevel = "Intermediate"
	ExperienceLevelPro          ExperienceLevel = "Pro"
)

// WorkType represents the work type enum
type WorkType string

const (
	WorkTypePersonal WorkType = "Personal"
	WorkTypeFirm     WorkType = "Firm"
)

// Survey represents a survey entry
type Survey struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	Email           string          `json:"email" db:"email"`
	NatureOfWork    NatureOfWork    `json:"nature_of_work" db:"nature_of_work"`
	TradingType     TradingType     `json:"trading_type" db:"trading_type"`
	ExperienceLevel ExperienceLevel `json:"experience_level" db:"experience_level"`
	WorkType        WorkType        `json:"work_type" db:"work_type"`
	Age             int             `json:"age" db:"age"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

// NewSurvey creates a new survey entry
func NewSurvey(email string, natureOfWork NatureOfWork, tradingType TradingType, experienceLevel ExperienceLevel, workType WorkType, age int) *Survey {
	now := time.Now()
	return &Survey{
		Email:           email,
		NatureOfWork:    natureOfWork,
		TradingType:     tradingType,
		ExperienceLevel: experienceLevel,
		WorkType:        workType,
		Age:             age,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
