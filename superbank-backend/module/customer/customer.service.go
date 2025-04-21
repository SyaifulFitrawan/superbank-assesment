package customer

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"
	"math"
)

type customerServiceImpl struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerServiceImpl{repo: repo}
}

func (s *customerServiceImpl) Create(ctx context.Context, input CustomerCreateRequest) (*model.Customer, error) {
	number := utils.GenerateBankAccountNumber()

	customer := &model.Customer{
		Name:          input.Name,
		Phone:         input.Phone,
		Address:       input.Address,
		ParentName:    input.ParentName,
		AccountNumber: number,
		AccountBranch: input.AccountBranch,
		AccountType:   input.AccountType,
	}

	err := s.repo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerServiceImpl) List(ctx context.Context, page, limit int, search string) ([]model.Customer, *utils.Paginator, error) {
	offset := (page - 1) * limit
	customers, total, err := s.repo.List(ctx, limit, offset, search)
	if err != nil {
		return nil, nil, err
	}

	pageCount := int(math.Ceil(float64(total) / float64(limit)))
	hasPrev := page > 1
	hasNext := page < pageCount

	var prevPage, nextPage *int
	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	if hasNext {
		n := page + 1
		nextPage = &n
	}

	paginator := &utils.Paginator{
		ItemCount:   int(total),
		Limit:       limit,
		PageCount:   pageCount,
		Page:        page,
		HasPrevPage: hasPrev,
		HasNextPage: hasNext,
		PrevPage:    prevPage,
		NextPage:    nextPage,
	}

	return customers, paginator, nil
}

func (s *customerServiceImpl) Detail(ctx context.Context, id string) (*CustomerDetailResponse, error) {
	customer, err := s.repo.Detail(ctx, id)
	if err != nil {
		return nil, err
	}

	result := &CustomerDetailResponse{
		Customer: *customer,
		Deposits: customer.Deposits,
		Pockets:  customer.Pockets,
	}

	return result, nil
}

func (s *customerServiceImpl) Update(ctx context.Context, id string, input CustomerUpdateRequest) error {
	customer := &model.Customer{}

	if input.Name != "" {
		customer.Name = input.Name
	}
	if input.Phone != "" {
		customer.Phone = input.Phone
	}
	if input.Address != "" {
		customer.Address = input.Address
	}
	if input.ParentName != "" {
		customer.ParentName = input.ParentName
	}
	if input.AccountBranch != "" {
		customer.AccountBranch = input.AccountBranch
	}
	if input.AccountType != "" {
		customer.AccountType = input.AccountType
	}

	if err := s.repo.Update(ctx, id, customer); err != nil {
		return err
	}

	return nil
}

func (s *customerServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
