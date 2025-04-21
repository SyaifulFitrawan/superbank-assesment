package pocket

import (
	"bank-backend/model"
	"bank-backend/module/customer"
	"context"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type pocketServiceImpl struct {
	pocketRepository   PocketRepository
	customerRepository customer.CustomerRepository
}

func NewPocketService(
	pocketRepository PocketRepository,
	customerRepository customer.CustomerRepository,
) PocketService {
	return &pocketServiceImpl{
		pocketRepository:   pocketRepository,
		customerRepository: customerRepository,
	}
}

func (s *pocketServiceImpl) Create(ctx context.Context, input CreatePocketRequest) (*model.Pocket, error) {
	payload := &model.Pocket{}

	payload.CustomerID = uuid.FromStringOrNil(input.CustomerID)
	payload.Name = input.Name

	if input.TargetAmount != nil {
		payload.TargetAmount = input.TargetAmount
	}

	if input.TargetDate != nil {
		targetDate, err := time.Parse("2006-1-2", *input.TargetDate)
		if err != nil {
			return nil, errors.New("invalid target date format")
		}

		payload.TargetDate = &targetDate
	}

	err := s.pocketRepository.Create(ctx, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s *pocketServiceImpl) TopUp(ctx context.Context, id string, input TopUpOrWithdrawPocketRequest) error {
	pocket, err := s.pocketRepository.Detail(ctx, id)
	if err != nil {
		return err
	}

	if !pocket.IsActive {
		message := fmt.Sprintf("Pocket %s is inactive", pocket.Name)
		return errors.New(message)
	}

	customer, err := s.customerRepository.Detail(ctx, pocket.CustomerID.String())
	if err != nil {
		return err
	}

	if customer.Balance < input.Amount {
		message := fmt.Sprintf("%s not have enough balance", customer.Name)
		return errors.New(message)
	}

	pocketPayload := model.Pocket{
		Balance: input.Amount + pocket.Balance,
	}

	customerPayload := model.Customer{
		Balance: customer.Balance - input.Amount,
	}

	err = s.pocketRepository.Update(ctx, id, &pocketPayload)
	if err != nil {
		return err
	}

	err = s.customerRepository.Update(ctx, customer.ID.String(), &customerPayload)
	if err != nil {
		return err
	}

	return nil
}

func (s *pocketServiceImpl) Withdrawn(ctx context.Context, id string, input TopUpOrWithdrawPocketRequest) error {
	pocket, err := s.pocketRepository.Detail(ctx, id)
	if err != nil {
		return err
	}

	if !pocket.IsActive {
		message := fmt.Sprintf("Pocket %s is inactive", pocket.Name)
		return errors.New(message)
	}

	customer, err := s.customerRepository.Detail(ctx, pocket.CustomerID.String())
	if err != nil {
		return err
	}

	if pocket.Balance < input.Amount {
		message := fmt.Sprintf("%s pocket not have enough balance", customer.Name)
		return errors.New(message)
	}

	pocketPayload := model.Pocket{
		Balance: pocket.Balance - input.Amount,
	}

	customerPayload := model.Customer{
		Balance: customer.Balance + input.Amount,
	}

	err = s.pocketRepository.Update(ctx, id, &pocketPayload)
	if err != nil {
		return err
	}

	err = s.customerRepository.Update(ctx, customer.ID.String(), &customerPayload)
	if err != nil {
		return err
	}

	return nil
}

func (s *pocketServiceImpl) Deactivated(ctx context.Context, id string) error {
	pocket, err := s.pocketRepository.Detail(ctx, id)
	if err != nil {
		return err
	}

	if !pocket.IsActive {
		message := fmt.Sprintf("Pocket %s has been inactive", pocket.Name)
		return errors.New(message)
	}

	if pocket.Balance != 0 {
		customer, err := s.customerRepository.Detail(ctx, pocket.CustomerID.String())
		if err != nil {
			return err
		}

		err = s.customerRepository.Update(ctx, customer.ID.String(), &model.Customer{
			Balance: customer.Balance + pocket.Balance,
		})
		if err != nil {
			return err
		}
	}

	err = s.pocketRepository.Deactivated(ctx, pocket.ID.String(), map[string]interface{}{
		"balance":   0,
		"is_active": false,
	})
	if err != nil {
		return err
	}

	return nil
}
