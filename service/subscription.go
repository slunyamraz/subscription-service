package service

import (
	"context"
	"database/sql"
	"s/domain"
	"s/repository"
)

type SubscriptionSerivce struct {
	db         *sql.DB
	subRepo    *repository.SubscriptionRepo
	outboxRepo *repository.OutboxRepo
}

func NewSubscriptionSerivce(db *sql.DB, subrepo *repository.SubscriptionRepo, outboxrepo *repository.OutboxRepo) *SubscriptionSerivce {
	return &SubscriptionSerivce{
		db:         db,
		subRepo:    subrepo,
		outboxRepo: outboxrepo,
	}
}

func (s *SubscriptionSerivce) GetByUserID(userID string) ([]domain.Subscription, error) {
	subs, err := s.subRepo.GetByUserID(userID)
	return subs, err
}

func (s *SubscriptionSerivce) GetListRepo() ([]domain.Subscription, error) {
	subs, err := s.subRepo.GetListRepo()
	return subs, err
}

func (s *SubscriptionSerivce) GetTotalPrice(user domain.UserTR) (string, error) {
	return s.subRepo.GetTotalPrice(user)
}

func (s *SubscriptionSerivce) Create(ctx context.Context, req domain.CreateSubscriptionRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	id, err := s.subRepo.Create(tx, req)
	if err != nil {
		return err
	}

	outboxEvent := domain.OutboxEvent{
		AggregateID:   id,
		AggregateType: domain.AggregateTypeSubscription,
		EventType:     domain.EventTypeSubscriptionCreated,
		Payload:       req,
	}
	err = s.outboxRepo.CreateEvent(tx, outboxEvent)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SubscriptionSerivce) Delete(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	outboxEvent := domain.OutboxEvent{
		AggregateID:   id,
		AggregateType: domain.AggregateTypeSubscription,
		EventType:     domain.EventTypeSubscriptionDeleted,
		Payload:       map[string]string{"id": id},
	}
	err = s.outboxRepo.CreateEvent(tx, outboxEvent)
	if err != nil {
		return err
	}

	err = s.subRepo.DeleteByID(tx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SubscriptionSerivce) UpdateByID(ctx context.Context, req domain.UpdateSubscriptionRequest) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	outboxEvent := domain.OutboxEvent{
		AggregateID:   req.ID,
		AggregateType: domain.AggregateTypeSubscription,
		EventType:     domain.EventTypeSubscriptionUpdated,
		Payload:       map[string]string{"id": req.ID},
	}

	err = s.outboxRepo.CreateEvent(tx, outboxEvent)
	if err != nil {
		return err
	}

	err = s.subRepo.UpdateByID(tx, req)
	if err != nil {
		return err
	}

	return tx.Commit()
}
