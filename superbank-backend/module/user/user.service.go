package user

import (
	"bank-backend/model"
	"bank-backend/utils"
	"context"
	"errors"
	"math"

	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

var HashPasswordFunc = bcrypt.GenerateFromPassword

func (s *userServiceImpl) Create(ctx context.Context, input *model.User) (*model.User, error) {
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("email and password are required")
	}

	hashedPassword, err := HashPasswordFunc([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	input.Password = string(hashedPassword)

	err = s.repo.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (s *userServiceImpl) List(ctx context.Context, page, limit int, search string) ([]model.User, *utils.Paginator, error) {
	offset := (page - 1) * limit
	users, total, err := s.repo.List(ctx, limit, offset, search)
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

	return users, paginator, nil
}

func (s *userServiceImpl) Detail(ctx context.Context, id string) (*model.User, error) {
	return s.repo.Detail(ctx, id)
}

func (s *userServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
