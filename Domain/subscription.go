package Domain

type Subscription struct {
	ID           string `json:"id"`
	Service_name string `json:"service_name"`
	Price        int    `json:"price"`
	UserID       string `json:"user_id"`
	StartDate    string `json:"start_date"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
}

type UpdateSubscriptionRequest struct {
	ID        string
	Price     int    `json:"price,omitempty"`
	StartDate string `json:"start_date,omitempty"`
}

type UserTR struct {
	UserID      string
	ServiceName string
	StartDate   string
}
