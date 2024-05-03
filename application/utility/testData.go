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
		Name:        "Насос CMP-151SG1 (титан)",
		Brand:       "Beto",
		Description: "Насос напольный ручной с манометром.",
		Image:       product1_defaultImage,
		Price:       54.00,
	}

	product2_defaultImage, err := GetDefaultImageBytes("default_images/velo.png", "png")
	if err != nil {
		fmt.Println("Unable to get velo.png")
		return
	}

	product2 := &entities.Product{
		Name:        "Велосипед Avalanche Comp 29 (чёрный)",
		Brand:       "GT",
		Description: "Легендарная модель американского бренда GT, которая выпускается уже много лет подряд, Avalanche Comp в последнем рестайлинге перетерпела ряд изменений, которые актуализировали велосипед в соответствии с нынешними тенденциями велоиндустрии. Рама байка сохранила фирменную технологию Triple Triangle™, которая обеспечивает велосипеду не только узнаваемость, но и высокую торсионную жесткость. Также по-прежнему трубы рамы выполнены с применением технологий баттинга и гидроформинга. Благодаря им рама легкая и жесткая. Из нового у велосипеда появился конусный рулевой стакан. Он существенно повышает торсионную жесткость узла, да и нижний подшипник рулевой колонки, который испытывает максимум нагрузки, стал крепче за счет большего диаметра.",
		Image:       product2_defaultImage,
		Price:       680.00,
	}

	category1 := &entities.Category{
		Name: "Насосы",
	}

	category2 := &entities.Category{
		Name: "Велосипеды",
	}

	//TODO: Adding new products via uc here
	//s.AddProductTr(product1, category1)
	//s.AddProductTr(product2, category2)
	uc.Begin()
	product1ID, err := uc.Products().Create(product1)
	if err != nil {
		uc.Rollback()
		fmt.Println("Unable to add product: ", err.Error())
		return
	}
	uc.Categories().Create(category1)
	uc.Products().AssignCategoryByName(product1ID, category1.Name)

	product2ID, err := uc.Products().Create(product2)
	if err != nil {
		uc.Rollback()
		fmt.Println("Unable to add product: ", err.Error())
		return
	}
	uc.Categories().Create(category2)
	uc.Products().AssignCategoryByName(product2ID, category2.Name)

	uc.Commit()
}
