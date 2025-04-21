package deposit

import (
	"bank-backend/model"
	"bank-backend/module/customer"
	"bank-backend/utils"
	"context"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

type depositServiceImpl struct {
	depositRepository  DepositRepository
	customerRepository customer.CustomerRepository
}

func NewDepositService(
	depositRepository DepositRepository,
	customerRepository customer.CustomerRepository,
) DepositService {
	return &depositServiceImpl{
		depositRepository:  depositRepository,
		customerRepository: customerRepository,
	}
}

func (s *depositServiceImpl) Create(ctx context.Context, input CreateDepositRequest) (*model.Deposit, error) {
	customer, err := s.customerRepository.Detail(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}

	if customer.Balance < input.Amount {
		message := fmt.Sprintf("%s not have enough balance", customer.Name)
		return nil, errors.New(message)
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	update := &model.Customer{
		Balance: customer.Balance - input.Amount,
	}

	err = s.customerRepository.Update(ctx, input.CustomerID, update)
	if err != nil {
		return nil, err
	}

	maturityDate := startDate.AddDate(0, input.TermMonths, 0)

	deposit := &model.Deposit{
		CustomerID:   uuid.FromStringOrNil(input.CustomerID),
		Amount:       input.Amount,
		InterestRate: input.InterestRate,
		TermMonths:   input.TermMonths,
		StartDate:    startDate,
		MaturityDate: maturityDate,
		IsWithdrawn:  false,
		Note:         input.Note,
	}

	err = s.depositRepository.Create(ctx, deposit)
	if err != nil {
		return nil, err
	}

	return deposit, nil
}

func (s *depositServiceImpl) ProcessMatureDeposits(ctx context.Context) error {
	logger := utils.NewLogger()
	deposits, err := s.depositRepository.FindMatureUnwithdraw(ctx)
	if err != nil {
		return nil
	}

	for _, deposit := range deposits {
		interest := (deposit.Amount * deposit.InterestRate * float64(deposit.TermMonths)) / 12
		total := deposit.Amount + interest

		if err := s.customerRepository.AddBalance(ctx, deposit.CustomerID.String(), total); err != nil {
			message := fmt.Sprintf("Failed to update balance for %s", deposit.Customer.Name)
			logger.Error(message, err.Error())
			continue
		}

		deposit.IsWithdrawn = true
		if err := s.depositRepository.Update(ctx, deposit.ID.String(), &deposit); err != nil {
			message := fmt.Sprintf("Failed to mark deposit %s", deposit.ID)
			logger.Error(message, err.Error())
		}
	}

	return nil
}
