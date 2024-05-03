package domain

import (
	"gorm.io/gorm"
	"kurs-server/domain/repos"
)

type UseCases struct {
	categories *repos.CategoryRepo
	products   *repos.ProductRepo
	users      *repos.UserRepo
}

func NewUseCases(storage *gorm.DB) *UseCases {
	return &UseCases{
		categories: &repos.CategoryRepo{
			Storage: storage,
		},
		products: &repos.ProductRepo{
			Storage: storage,
		},
		users: &repos.UserRepo{
			Storage: storage,
		},
	}
}

func (uc *UseCases) Categories() *repos.CategoryRepo {
	return uc.categories
}

func (uc *UseCases) Products() *repos.ProductRepo {
	return uc.products
}

func (uc *UseCases) Users() *repos.UserRepo {
	return uc.users
}

func (uc *UseCases) Begin() error {
	return nil
}

func (uc *UseCases) Commit() error { return nil }

func (uc *UseCases) Rollback() error { return nil }
