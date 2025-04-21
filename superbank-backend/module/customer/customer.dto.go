package customer

import "bank-backend/model"

type CustomerCreateRequest struct {
	Name          string `json:"name" validate:"required,notblank"`
	Phone         string `json:"phone" validate:"required,notblank"`
	Address       string `json:"address" validate:"required,notblank"`
	ParentName    string `json:"parent_name" validate:"required,notblank"`
	AccountBranch string `json:"account_branch" validate:"required,notblank"`
	AccountType   string `json:"account_type" validate:"required,notblank"`
}

type CustomerUpdateRequest struct {
	Name          string `json:"name" validate:"omitempty"`
	Phone         string `json:"phone" validate:"omitempty"`
	Address       string `json:"address" validate:"omitempty"`
	ParentName    string `json:"parent_name" validate:"omitempty"`
	AccountBranch string `json:"account_branch" validate:"omitempty"`
	AccountType   string `json:"account_type" validate:"omitempty"`
}

type CustomerDetailResponse struct {
	model.Customer
	Deposits []model.Deposit `json:"deposits"`
	Pockets  []model.Pocket  `json:"pockets"`
}
