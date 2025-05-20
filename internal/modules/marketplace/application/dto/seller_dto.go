// filepath: /Users/amit/workspace/sutram/internal/modules/marketplace/application/dto/seller_dto.go
package dto

// SellerAnalyticsResponse represents the analytics data for a seller
type SellerAnalyticsResponse struct {
	TotalRevenue        float64 `json:"total_revenue"`
	TotalSubscriptions  int     `json:"total_subscriptions"`
	ActiveSubscriptions int     `json:"active_subscriptions"`
	TotalListings       int     `json:"total_listings"`
	ActiveListings      int     `json:"active_listings"`
	AverageRating       float32 `json:"average_rating"`
	TotalViews          int     `json:"total_views"`
	// Revenue breakdown by time period
	RevenueByMonth map[string]float64 `json:"revenue_by_month"`
}
