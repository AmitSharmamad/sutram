package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"algo/internal/modules/strategy/application/dto"
	"algo/internal/modules/strategy/application/service"
	strategyHttp "algo/internal/modules/strategy/interfaces/http"
	userHttp "algo/internal/modules/user/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStrategyService is a mock implementation of the StrategyService interface
type MockStrategyService struct {
	mock.Mock
}

func (m *MockStrategyService) CreateStrategy(ctx context.Context, req *dto.CreateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StrategyResponse), args.Error(1)
}

func (m *MockStrategyService) GetStrategyByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StrategyResponse), args.Error(1)
}

func (m *MockStrategyService) GetUserStrategies(ctx context.Context, userID uuid.UUID) ([]*dto.StrategyResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.StrategyResponse), args.Error(1)
}

func (m *MockStrategyService) UpdateStrategy(ctx context.Context, req *dto.UpdateStrategyRequest, userID uuid.UUID) (*dto.StrategyResponse, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StrategyResponse), args.Error(1)
}

func (m *MockStrategyService) DeleteStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockStrategyService) ActivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockStrategyService) DeactivateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockStrategyService) DuplicateStrategy(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*dto.StrategyResponse, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.StrategyResponse), args.Error(1)
}

func (m *MockStrategyService) SearchStrategies(ctx context.Context, filter *dto.StrategyFilter, userID uuid.UUID) ([]*dto.StrategyResponse, error) {
	args := m.Called(ctx, filter, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dto.StrategyResponse), args.Error(1)
}

// MockAuthService is a mock implementation of the auth service
type MockAuthService struct{}

func (m *MockAuthService) Authenticate(ctx context.Context, token string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

// MockAuthServiceAdapter is a mock adapter for the auth service
type MockAuthServiceAdapter struct{}

func (m *MockAuthServiceAdapter) Authenticate(c *gin.Context, token string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}

func (m *MockAuthServiceAdapter) AuthenticateRequest(c *gin.Context) {
	userID := uuid.New()
	c.Set("userID", userID)
}

// MockUserController is a mock implementation of the UserController
type MockUserController struct{}

func (m *MockUserController) GetAuthService() interface{} {
	return &MockAuthService{}
}

// NewAuthServiceAdapter returns a mock auth service adapter
func NewAuthServiceAdapter(authService interface{}) *MockAuthServiceAdapter {
	return &MockAuthServiceAdapter{}
}

// Mock the AuthServiceAdapter function
func init() {
	// Save the original function
	originalNewAuthServiceAdapter := userHttp.NewAuthServiceAdapter
	
	// Override for testing
	userHttp.NewAuthServiceAdapter = func(authService interface{}) interface{} {
		return &MockAuthServiceAdapter{}
	}
	
	// Restore the original in a cleanup function (may not be called in tests)
	// This is just to avoid affecting other packages if they run in the same process
	if false { // never executed but keeps the reference to originalNewAuthServiceAdapter
		userHttp.NewAuthServiceAdapter = originalNewAuthServiceAdapter
	}
}

// TestNewStrategyHandler tests the creation of a new strategy handler
func TestNewStrategyHandler(t *testing.T) {
	// Arrange
	mockService := new(MockStrategyService)

	// Act
	handler := strategyHttp.NewStrategyHandler(mockService)

	// Assert
	assert.NotNil(t, handler)
}

// TestRegisterRoutes tests the route registration
func TestRegisterRoutes(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	handler := strategyHttp.NewStrategyHandler(mockService)
	router := gin.New()
	group := router.Group("/api")
	userController := &MockUserController{}

	// Act & Assert - Should not panic
	assert.NotPanics(t, func() {
		handler.RegisterRoutes(group, userController)
	})
}

// TestCreateStrategy tests the CreateStrategy handler
func TestCreateStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockReq := &dto.CreateStrategyRequest{
			Name:             "Test Strategy",
			Description:      "Test Description",
			PositionType:     "LONG",
			Instruments:      []int32{1, 2, 3},
			ChartType:        "CANDLESTICK",
			ChartTimeframe:   "1h",
			TradeType:        "MARKET",
			Qty:              10,
			EntryExpression:  "close > open",
			ExitExpression:   "close < open",
			TargetProfit:     5.0,
			TargetProfitType: "PERCENTAGE",
			StopLoss:         2.0,
			StopLossType:     "PERCENTAGE",
			Capital:          10000.0,
		}

		mockResp := &dto.StrategyResponse{
			ID:               uuid.New(),
			UserID:           userID,
			Name:             mockReq.Name,
			Description:      mockReq.Description,
			PositionType:     mockReq.PositionType,
			ChartType:        mockReq.ChartType,
			ChartTimeframe:   mockReq.ChartTimeframe,
			TradeType:        mockReq.TradeType,
			EntryExpression:  mockReq.EntryExpression,
			ExitExpression:   mockReq.ExitExpression,
			TargetProfit:     mockReq.TargetProfit,
			TargetProfitType: mockReq.TargetProfitType,
			StopLoss:         mockReq.StopLoss,
			StopLossType:     mockReq.StopLossType,
			Capital:          mockReq.Capital,
			IsActive:         true,
		}

		mockService.On("CreateStrategy", mock.Anything, mock.MatchedBy(func(req *dto.CreateStrategyRequest) bool {
			return req.Name == mockReq.Name
		}), userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.CreateStrategy(c)
		})

		// Create request
		reqBody, _ := json.Marshal(mockReq)
		req, _ := http.NewRequest(http.MethodPost, "/strategy", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy", handler.CreateStrategy)

		// Create request
		reqBody, _ := json.Marshal(map[string]string{"name": "Test"})
		req, _ := http.NewRequest(http.MethodPost, "/strategy", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		// Arrange
		mockService.On("CreateStrategy", mock.Anything, mock.Anything, userID).
			Return(nil, service.ErrInvalidRequest)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.CreateStrategy(c)
		})

		// Create invalid request (empty body)
		reqBody, _ := json.Marshal(map[string]string{})
		req, _ := http.NewRequest(http.MethodPost, "/strategy", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

// TestGetStrategyByID tests the GetStrategyByID handler
func TestGetStrategyByID(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockResp := &dto.StrategyResponse{
			ID:               strategyID,
			UserID:           userID,
			Name:             "Test Strategy",
			Description:      "Test Description",
			PositionType:     "LONG",
			ChartType:        "CANDLESTICK",
			ChartTimeframe:   "1h",
			TradeType:        "MARKET",
			EntryExpression:  "close > open",
			ExitExpression:   "close < open",
			TargetProfit:     5.0,
			TargetProfitType: "PERCENTAGE",
			StopLoss:         2.0,
			StopLossType:     "PERCENTAGE",
			Capital:          10000.0,
			IsActive:         true,
		}

		mockService.On("GetStrategyByID", mock.Anything, strategyID, userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetStrategyByID(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy/:id", handler.GetStrategyByID)

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetStrategyByID(c)
		})

		// Create request with invalid ID
		req, _ := http.NewRequest(http.MethodGet, "/strategy/invalid-id", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockService.On("GetStrategyByID", mock.Anything, strategyID, userID).
			Return(nil, service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetStrategyByID(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestGetUserStrategies tests the GetUserStrategies handler
func TestGetUserStrategies(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockResp := []*dto.StrategyResponse{
			{
				ID:               uuid.New(),
				UserID:           userID,
				Name:             "Test Strategy 1",
				Description:      "Test Description 1",
				PositionType:     "LONG",
				ChartType:        "CANDLESTICK",
				ChartTimeframe:   "1h",
				TradeType:        "MARKET",
				EntryExpression:  "close > open",
				ExitExpression:   "close < open",
				TargetProfit:     5.0,
				TargetProfitType: "PERCENTAGE",
				StopLoss:         2.0,
				StopLossType:     "PERCENTAGE",
				Capital:          10000.0,
				IsActive:         true,
			},
			{
				ID:               uuid.New(),
				UserID:           userID,
				Name:             "Test Strategy 2",
				Description:      "Test Description 2",
				PositionType:     "SHORT",
				ChartType:        "HEIKINASHI",
				ChartTimeframe:   "5m",
				TradeType:        "LIMIT",
				EntryExpression:  "macd > signal",
				ExitExpression:   "macd < signal",
				TargetProfit:     10.0,
				TargetProfitType: "ABSOLUTE",
				StopLoss:         5.0,
				StopLossType:     "ABSOLUTE",
				Capital:          20000.0,
				IsActive:         false,
			},
		}

		mockService.On("GetUserStrategies", mock.Anything, userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetUserStrategies(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Empty Result", func(t *testing.T) {
		// Arrange
		mockService.On("GetUserStrategies", mock.Anything, userID).Return([]*dto.StrategyResponse{}, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetUserStrategies(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)

		// Check response body - should be an empty array, not null
		var respObj map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &respObj)
		assert.NoError(t, err)
		assert.IsType(t, []interface{}{}, respObj["data"])
		assert.Len(t, respObj["data"], 0)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Nil Result", func(t *testing.T) {
		// Arrange - service returns nil instead of empty slice
		mockService.On("GetUserStrategies", mock.Anything, userID).Return(nil, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetUserStrategies(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)

		// Check response body - should be an empty array, not null
		var respObj map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &respObj)
		assert.NoError(t, err)
		assert.IsType(t, []interface{}{}, respObj["data"])
		assert.Len(t, respObj["data"], 0)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy", handler.GetUserStrategies)

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Arrange
		mockService.On("GetUserStrategies", mock.Anything, userID).
			Return(nil, errors.New("database error"))

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.GET("/strategy", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.GetUserStrategies(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/strategy", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestUpdateStrategy tests the UpdateStrategy handler
func TestUpdateStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockReq := &dto.UpdateStrategyRequest{
			ID:               strategyID,
			Name:             "Updated Strategy",
			Description:      "Updated Description",
			PositionType:     "SHORT",
			Instruments:      []int32{4, 5, 6},
			ChartType:        "HEIKINASHI",
			ChartTimeframe:   "5m",
			TradeType:        "LIMIT",
			Qty:              20,
			EntryExpression:  "macd > signal",
			ExitExpression:   "macd < signal",
			TargetProfit:     10.0,
			TargetProfitType: "ABSOLUTE",
			StopLoss:         5.0,
			StopLossType:     "ABSOLUTE",
			Capital:          20000.0,
		}

		mockResp := &dto.StrategyResponse{
			ID:               strategyID,
			UserID:           userID,
			Name:             mockReq.Name,
			Description:      mockReq.Description,
			PositionType:     mockReq.PositionType,
			ChartType:        mockReq.ChartType,
			ChartTimeframe:   mockReq.ChartTimeframe,
			TradeType:        mockReq.TradeType,
			EntryExpression:  mockReq.EntryExpression,
			ExitExpression:   mockReq.ExitExpression,
			TargetProfit:     mockReq.TargetProfit,
			TargetProfitType: mockReq.TargetProfitType,
			StopLoss:         mockReq.StopLoss,
			StopLossType:     mockReq.StopLossType,
			Capital:          mockReq.Capital,
			IsActive:         true,
		}

		mockService.On("UpdateStrategy", mock.Anything, mock.MatchedBy(func(req *dto.UpdateStrategyRequest) bool {
			return req.ID == strategyID && req.Name == mockReq.Name
		}), userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.PUT("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.UpdateStrategy(c)
		})

		// Create request
		reqBody, _ := json.Marshal(mockReq)
		req, _ := http.NewRequest(http.MethodPut, "/strategy/"+strategyID.String(), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.PUT("/strategy/:id", handler.UpdateStrategy)

		// Create request
		reqBody, _ := json.Marshal(map[string]string{"name": "Test"})
		req, _ := http.NewRequest(http.MethodPut, "/strategy/"+strategyID.String(), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.PUT("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.UpdateStrategy(c)
		})

		// Create request with invalid ID
		reqBody, _ := json.Marshal(map[string]string{"name": "Test"})
		req, _ := http.NewRequest(http.MethodPut, "/strategy/invalid-id", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockReq := &dto.UpdateStrategyRequest{
			ID:               strategyID,
			Name:             "Updated Strategy",
			Description:      "Updated Description",
			PositionType:     "SHORT",
			Instruments:      []int32{4, 5, 6},
			ChartType:        "HEIKINASHI",
			ChartTimeframe:   "5m",
			TradeType:        "LIMIT",
			Qty:              20,
			EntryExpression:  "macd > signal",
			ExitExpression:   "macd < signal",
			TargetProfit:     10.0,
			TargetProfitType: "ABSOLUTE",
			StopLoss:         5.0,
			StopLossType:     "ABSOLUTE",
			Capital:          20000.0,
		}

		mockService.On("UpdateStrategy", mock.Anything, mock.MatchedBy(func(req *dto.UpdateStrategyRequest) bool {
			return req.ID == strategyID
		}), userID).Return(nil, service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.PUT("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.UpdateStrategy(c)
		})

		// Create request
		reqBody, _ := json.Marshal(mockReq)
		req, _ := http.NewRequest(http.MethodPut, "/strategy/"+strategyID.String(), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestDeleteStrategy tests the DeleteStrategy handler
func TestDeleteStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockService.On("DeleteStrategy", mock.Anything, strategyID, userID).Return(nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.DELETE("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeleteStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodDelete, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.DELETE("/strategy/:id", handler.DeleteStrategy)

		// Create request
		req, _ := http.NewRequest(http.MethodDelete, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.DELETE("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeleteStrategy(c)
		})

		// Create request with invalid ID
		req, _ := http.NewRequest(http.MethodDelete, "/strategy/invalid-id", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockService.On("DeleteStrategy", mock.Anything, strategyID, userID).
			Return(service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.DELETE("/strategy/:id", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeleteStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodDelete, "/strategy/"+strategyID.String(), nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestActivateStrategy tests the ActivateStrategy handler
func TestActivateStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockService.On("ActivateStrategy", mock.Anything, strategyID, userID).Return(nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/activate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.ActivateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/activate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/activate", handler.ActivateStrategy)

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/activate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/activate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.ActivateStrategy(c)
		})

		// Create request with invalid ID
		req, _ := http.NewRequest(http.MethodPost, "/strategy/invalid-id/activate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockService.On("ActivateStrategy", mock.Anything, strategyID, userID).
			Return(service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/activate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.ActivateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/activate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestDeactivateStrategy tests the DeactivateStrategy handler
func TestDeactivateStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockService.On("DeactivateStrategy", mock.Anything, strategyID, userID).Return(nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/deactivate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeactivateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/deactivate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/deactivate", handler.DeactivateStrategy)

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/deactivate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/deactivate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeactivateStrategy(c)
		})

		// Create request with invalid ID
		req, _ := http.NewRequest(http.MethodPost, "/strategy/invalid-id/deactivate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockService.On("DeactivateStrategy", mock.Anything, strategyID, userID).
			Return(service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/deactivate", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DeactivateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/deactivate", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestDuplicateStrategy tests the DuplicateStrategy handler
func TestDuplicateStrategy(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()
	strategyID := uuid.New()
	newStrategyID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		mockResp := &dto.StrategyResponse{
			ID:               newStrategyID,
			UserID:           userID,
			Name:             "Test Strategy (Copy)",
			Description:      "Test Description",
			PositionType:     "LONG",
			ChartType:        "CANDLESTICK",
			ChartTimeframe:   "1h",
			TradeType:        "MARKET",
			EntryExpression:  "close > open",
			ExitExpression:   "close < open",
			TargetProfit:     5.0,
			TargetProfitType: "PERCENTAGE",
			StopLoss:         2.0,
			StopLossType:     "PERCENTAGE",
			Capital:          10000.0,
			IsActive:         true,
		}

		mockService.On("DuplicateStrategy", mock.Anything, strategyID, userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/clone", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DuplicateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/clone", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/clone", handler.DuplicateStrategy)

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/clone", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/clone", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DuplicateStrategy(c)
		})

		// Create request with invalid ID
		req, _ := http.NewRequest(http.MethodPost, "/strategy/invalid-id/clone", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Arrange
		mockService.On("DuplicateStrategy", mock.Anything, strategyID, userID).
			Return(nil, service.ErrStrategyNotFound)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/:id/clone", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.DuplicateStrategy(c)
		})

		// Create request
		req, _ := http.NewRequest(http.MethodPost, "/strategy/"+strategyID.String()+"/clone", nil)
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

// TestSearchStrategies tests the SearchStrategies handler
func TestSearchStrategies(t *testing.T) {
	// Setup common test elements
	gin.SetMode(gin.TestMode)
	mockService := new(MockStrategyService)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		// Arrange
		filter := &dto.StrategyFilter{
			Name:         "Test",
			PositionType: "LONG",
			IsActive:     true,
		}

		mockResp := []*dto.StrategyResponse{
			{
				ID:               uuid.New(),
				UserID:           userID,
				Name:             "Test Strategy 1",
				Description:      "Test Description 1",
				PositionType:     "LONG",
				ChartType:        "CANDLESTICK",
				ChartTimeframe:   "1h",
				TradeType:        "MARKET",
				EntryExpression:  "close > open",
				ExitExpression:   "close < open",
				TargetProfit:     5.0,
				TargetProfitType: "PERCENTAGE",
				StopLoss:         2.0,
				StopLossType:     "PERCENTAGE",
				Capital:          10000.0,
				IsActive:         true,
			},
		}

		mockService.On("SearchStrategies", mock.Anything, mock.MatchedBy(func(f *dto.StrategyFilter) bool {
			return f.Name == filter.Name && f.PositionType == filter.PositionType && f.IsActive == filter.IsActive
		}), userID).Return(mockResp, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.SearchStrategies(c)
		})

		// Create request
		reqBody, _ := json.Marshal(filter)
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Empty Result", func(t *testing.T) {
		// Arrange
		filter := &dto.StrategyFilter{
			Name: "NonExistent",
		}

		mockService.On("SearchStrategies", mock.Anything, mock.Anything, userID).
			Return([]*dto.StrategyResponse{}, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.SearchStrategies(c)
		})

		// Create request
		reqBody, _ := json.Marshal(filter)
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)

		// Check response body - should be an empty array, not null
		var respObj map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &respObj)
		assert.NoError(t, err)
		assert.IsType(t, []interface{}{}, respObj["data"])
		assert.Len(t, respObj["data"], 0)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Nil Result", func(t *testing.T) {
		// Arrange - service returns nil instead of empty slice
		filter := &dto.StrategyFilter{
			Name: "NonExistent",
		}

		mockService.On("SearchStrategies", mock.Anything, mock.Anything, userID).
			Return(nil, nil)

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.SearchStrategies(c)
		})

		// Create request
		reqBody, _ := json.Marshal(filter)
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusOK, resp.Code)

		// Check response body - should be an empty array, not null
		var respObj map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &respObj)
		assert.NoError(t, err)
		assert.IsType(t, []interface{}{}, respObj["data"])
		assert.Len(t, respObj["data"], 0)
		
		mockService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", handler.SearchStrategies)

		// Create request
		reqBody, _ := json.Marshal(map[string]string{"name": "Test"})
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		// Arrange
		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.SearchStrategies(c)
		})

		// Create request with invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		// Arrange
		filter := &dto.StrategyFilter{
			Name: "Test",
		}

		mockService.On("SearchStrategies", mock.Anything, mock.Anything, userID).
			Return(nil, errors.New("database error"))

		handler := strategyHttp.NewStrategyHandler(mockService)
		router := gin.New()
		router.POST("/strategy/search", func(c *gin.Context) {
			// Set userID in context (normally done by auth middleware)
			c.Set("userID", userID)
			handler.SearchStrategies(c)
		})

		// Create request
		reqBody, _ := json.Marshal(filter)
		req, _ := http.NewRequest(http.MethodPost, "/strategy/search", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		// Act
		router.ServeHTTP(resp, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}