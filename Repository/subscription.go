package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"s/domain"
	"strconv"
)

type SubscriptionRepo struct {
	db *sql.DB
}

func NewSubscriptionRepo(db *sql.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(tx *sql.Tx, req domain.CreateSubscriptionRequest) (string, error) {
	var id string
	err := tx.QueryRow(`INSERT INTO subscriptions (service_name, price, user_id, start_date) 
	          VALUES ($1, $2, $3, $4) RETURNING id`, req.ServiceName, req.Price, req.UserID, req.StartDate).Scan(&id)
	if err != nil {
		slog.Error("error in CreateSub", "error", err)
		return "", err
	}
	return id, nil
}

func (r *SubscriptionRepo) GetByUserID(userID string) ([]domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date 
	          FROM subscriptions WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(&sub.ID, &sub.Service_name, &sub.Price, &sub.UserID, &sub.StartDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SubscriptionRepo) GetListRepo() ([]domain.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date FROM subscriptions`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		err := rows.Scan(&sub.ID, &sub.Service_name, &sub.Price, &sub.UserID, &sub.StartDate)
		if err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SubscriptionRepo) DeleteByID(tx *sql.Tx, id string) error {
	status, err := tx.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		slog.Error("error in delete repo:", err)
		return err
	}

	rowsAffected, err := status.RowsAffected()
	if err != nil {
		slog.Error("error in delete repo:", err)
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found with id %s", id)
	}

	return nil
}

func (s *SubscriptionRepo) UpdateByID(tx *sql.Tx, upd domain.UpdateSubscriptionRequest) error {
	if upd.Price != 0 && upd.StartDate != "" {
		status, err := tx.Exec(
			"UPDATE subscriptions SET price = $1, start_date = $2 WHERE id = $3",
			upd.Price, upd.StartDate, upd.ID)
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}

		rowsAffected, err := status.RowsAffected()
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no subscription found with id %s", upd.ID)
		}
		return nil
	}

	if upd.Price == 0 {
		status, err := tx.Exec(
			"UPDATE subscriptions SET start_date = $1 WHERE id = $2",
			upd.StartDate, upd.ID)
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}

		rowsAffected, err := status.RowsAffected()
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no subscription found with id %s", upd.ID)
		}
		return nil
	}

	if upd.StartDate == "" {
		status, err := tx.Exec(
			"UPDATE subscriptions SET price = $1 WHERE id = $2",
			upd.Price, upd.ID)
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}

		rowsAffected, err := status.RowsAffected()
		if err != nil {
			slog.Error("error in delete repo:", err)
			return err
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no subscription found with id %s", upd.ID)
		}
		return nil
	}
	return nil
}

func (s *SubscriptionRepo) GetTotalPrice(user domain.UserTR) (string, error) {
	query := `SELECT SUM(price) FROM subscriptions WHERE user_id = $1
                                       and start_date = $2
                                       and service_name = $3`
	var price int
	err := s.db.QueryRow(query, user.UserID, user.StartDate,
		user.ServiceName).Scan(&price)
	if err != nil {
		slog.Error("error in get total price:", err)
		return "", err
	}

	return strconv.Itoa(price), nil
}
