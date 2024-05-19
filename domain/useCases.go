package domain

import (
	"kurs-server/domain/repos"

	"gorm.io/gorm"
)

type UseCases struct {
	categories *repos.CategoryRepo
	products   *repos.ProductRepo
	users      *repos.UserRepo
	details    *repos.DetailRepo
	ratings    *repos.RatingRepo
	reviews    *repos.ReviewRepo
	orders     *repos.OrderRepo
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
		details: &repos.DetailRepo{
			Storage: storage,
		},
		ratings: &repos.RatingRepo{
			Storage: storage,
		},
		reviews: &repos.ReviewRepo{
			Storage: storage,
		},
		orders: &repos.OrderRepo{
			Storage: storage,
		},
	}
}

func (uc *UseCases) Orders() *repos.OrderRepo {
	return uc.orders
}

func (uc *UseCases) Ratings() *repos.RatingRepo {
	return uc.ratings
}

func (uc *UseCases) Reviews() *repos.ReviewRepo {
	return uc.reviews
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

func (uc *UseCases) Details() *repos.DetailRepo {
	return uc.details
}

func (uc *UseCases) Begin() error {
	return nil
}

func (uc *UseCases) Commit() error { return nil }

func (uc *UseCases) Rollback() error { return nil }
