package utility

import (
	"fmt"
	"kurs-server/domain"
	"kurs-server/domain/entities"

	"golang.org/x/crypto/bcrypt"
)

func AddTestData(uc *domain.UseCases) {
	admin_password, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	admin_user := &entities.User{
		Role:     "Admin",
		Username: "admin",
		Email:    "admin@admin.admin",
		Password: string(admin_password),
	}
	uc.Users().Create(admin_user)

	product1_defaultImage, err := GetDefaultImageBytes("default_images/nasos.png", "png")
	if err != nil {
		fmt.Println("Unable to get nasos.png")
		fmt.Println(err.Error())
		return
	}

	product1 := &entities.Product{
		Name:        "TestName",
		Brand:       "TestBrand",
		Description: "TestDesc",
		Image:       product1_defaultImage,
		Price:       54.00,
	}

	//creating product
	uc.Products().Create(product1)

	//creating category
	category1 := &entities.Category{
		Name: "Oil",
	}
	uc.Categories().Create(category1)

	//creating detail
	uc.Details().Create("Volume", "Oil")

	//assigning category to product
	uc.Products().AssignCategoryByName(product1.ID, category1.Name)

	//assigning value to detail
	uc.Details().CreateValue(product1.ID, 1, "10l")

}
