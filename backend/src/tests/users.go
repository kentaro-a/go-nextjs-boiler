package tests

import (
	"app/model"
	"time"
)

func UsersFixture() []*model.User {
	return []*model.User{
		{
			ID:        1,
			Name:      "user1",
			Mail:      "user1@test.com",
			Password:  "222f971b48b727d385973f5c560448ec3dfb952e07555a08583fa764fdddc36e",
			StatusFlg: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "user2",
			Mail:      "user2@test.com",
			Password:  "db6154e27b71e0a402b36a8c84471fc0e29c27fa0e47d9c317d0f02b14cf9efb",
			StatusFlg: 0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Name:      "user3",
			Mail:      "user3@test.com",
			Password:  "222f971b48b727d385973f5c560448ec3dfb952e07555a08583fa764fdddc36e",
			StatusFlg: 1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}
