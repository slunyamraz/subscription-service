package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"s/domain"
	"s/service"
)

type SubscriptionHandler struct {
	serviceS *service.SubscriptionSerivce
}

func NewSubscriptionHandler(servicesub *service.SubscriptionSerivce) *SubscriptionHandler {
	return &SubscriptionHandler{serviceS: servicesub}
}

func RegisterRoutes(h *SubscriptionHandler) {
	http.HandleFunc("POST /post", h.PostHandler)
	http.HandleFunc("DELETE /subscriptions/{id}", h.DeleteHandler)
	http.HandleFunc("PUT /update/{id}", h.UpdateHandler)
	http.HandleFunc("GET /get", h.GetHandler)
	http.HandleFunc("GET /list", h.GetListHandler)
	http.HandleFunc("GET /gettotalprice", h.GetTotalPrice)
}

func ServeStart() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error(err.Error(), "err on handle lvl")
	}
}

// GetHandler godoc
// @Summary Get user subscriptions
// @Description Get all subscriptions for user
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} domain.Subscription
// @Router /get [get]
func (s *SubscriptionHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error("user_id is required")
		return
	}
	subs, err := s.serviceS.GetByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
}

// GetListHandler godoc
// @Summary List all subscriptions
// @Description Get all subscriptions in system
// @Tags subscriptions
// @Produce json
// @Success 200 {array} domain.Subscription
// @Router /list [get]
func (s *SubscriptionHandler) GetListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}
	subs, err := s.serviceS.GetListRepo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(subs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

}

// PostHandler godoc
// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce plain
// @Param subscription body domain.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {string} string "Subscription created"
// @Router /post [post]
func (s *SubscriptionHandler) PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}

	var subreq domain.CreateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&subreq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}

	err = s.serviceS.Create(r.Context(), subreq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteHandler godoc
// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce plain
// @Param id path string true "Subscription ID"
// @Success 200 {string} string "Subscription deleted"
// @Router /subscriptions/{id} [delete]
func (s *SubscriptionHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}

	id := r.PathValue("id")

	err := s.serviceS.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateHandler godoc
// @Summary Update subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce plain
// @Param id path string true "Subscription ID"
// @Param subscription body domain.UpdateSubscriptionRequest true "Update data"
// @Success 200 {string} string "Subscription updated"
// @Router /update/{id} [put]
func (s *SubscriptionHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}

	var updatereq domain.UpdateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&updatereq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}
	updatereq.ID = r.PathValue("id")
	err = s.serviceS.UpdateByID(r.Context(), updatereq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}
}

// GetTotalPrice godoc
// @Summary Get total price
// @Description Get total price with filters
// @Tags subscriptions
// @Produce plain
// @Param user_id query string true "User ID"
// @Param service_name query string true "Service name"
// @Param start_date query string true "Start date (MM-YYYY)"
// @Success 200 {string} string "Total price"
// @Router /gettotalprice [get]
func (s *SubscriptionHandler) GetTotalPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		slog.Error("not allowed request", "method", r.Method)
		return
	}
	var user domain.UserTR
	user.UserID = r.URL.Query().Get("user_id")
	user.ServiceName = r.URL.Query().Get("service_name")
	user.StartDate = r.URL.Query().Get("start_date")
	price, err := s.serviceS.GetTotalPrice(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "err on handle lvl")
		return
	}
	w.Write([]byte(price))
}
