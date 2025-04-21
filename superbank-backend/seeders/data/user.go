package seeder_data

import "bank-backend/model"

var DummyUser = []model.User{
	{
		Email:    "admin@example.com",
		Username: "admin",
		Password: "password",
	},
	{
		Email:    "employee@example.com",
		Username: "employee",
		Password: "password",
	},
}
