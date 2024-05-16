package utility

import (
	"fmt"
	"kurs-server/domain"
	"kurs-server/domain/entities"

	"golang.org/x/crypto/bcrypt"
)

func AddTestData(uc *domain.UseCases) {
	addTestUsers(uc)
	addTestCategories(uc)
	addTestDetails(uc)
	addTestProducts(uc)
}

func addTestUsers(uc *domain.UseCases) {
	admin_password, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	admin_user := &entities.User{
		Role:     "Admin",
		Username: "admin",
		Email:    "admin@admin.admin",
		Password: string(admin_password),
	}
	uc.Users().Create(admin_user)

	user_passwrod, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

	default_user := &entities.User{
		Role:     "User",
		Username: "user",
		Email:    "user@user.user",
		Password: string(user_passwrod),
	}
	uc.Users().Create(default_user)
}

func addTestCategories(uc *domain.UseCases) {
	category1 := &entities.Category{
		Name: "Велосипеды",
	}
	uc.Categories().Create(category1)

	category2 := &entities.Category{
		Name: "Замки велосипедные",
	}
	uc.Categories().Create(category2)

	category3 := &entities.Category{
		Name: "Насосы для велосипеда",
	}
	uc.Categories().Create(category3)

	category4 := &entities.Category{
		Name: "Щитки для велосипеда",
	}
	uc.Categories().Create(category4)
}

func addTestDetails(uc *domain.UseCases) {
	//1
	uc.Details().Create("Класс", "Велосипеды")         //1
	uc.Details().Create("Назначение", "Велосипеды")    //2
	uc.Details().Create("Цвет", "Велосипеды")          //3
	uc.Details().Create("Материал рамы", "Велосипеды") //4

	//2
	uc.Details().Create("Вес", "Замки велосипедные")       //5
	uc.Details().Create("Цвет", "Замки велосипедные")      //6
	uc.Details().Create("Тип замка", "Замки велосипедные") //7
	uc.Details().Create("Запирание", "Замки велосипедные") //8

	//3
	uc.Details().Create("Тип", "Насосы для велосипеда")  //9
	uc.Details().Create("Вес", "Насосы для велосипеда")  //10
	uc.Details().Create("Цвет", "Насосы для велосипеда") //11

	//4
	uc.Details().Create("Тип", "Щитки для велосипеда")           //12
	uc.Details().Create("Вес", "Щитки для велосипеда")           //13
	uc.Details().Create("Цвет", "Щитки для велосипеда")          //14
	uc.Details().Create("Диаметр колёс", "Щитки для велосипеда") //15
}

func addTestProducts(uc *domain.UseCases) {
	// ------------------velo1-------------------------------
	velo1_image, err := GetDefaultImageBytes("default_images/merida.png", "png")
	if err != nil {
		fmt.Println(err.Error())
	}
	velo1 := &entities.Product{
		Name:        "Велосипед Merida Big.Nine LTD (матовый зелёный/шампанское)",
		Brand:       "Merida",
		Description: "Назначение: мужской; Количество скоростей: 9; Манетки: Microshift Advent Trail SL-M9195-R 9-speed; Задний тормоз: дисковый гидравлический; Производитель: Merida",
		Image:       velo1_image,
		Price:       2312,
	}
	uc.Products().Create(velo1)

	uc.Products().AssignCategoryByName(1, "Велосипеды")

	uc.Details().CreateValue(1, 1, "Горный (МТВ)")
	uc.Details().CreateValue(2, 1, "Мужской")
	uc.Details().CreateValue(3, 1, "Зеленый")
	uc.Details().CreateValue(4, 1, "Алюминий")

	// ------------------velo2-------------------------------
	velo2_image, _ := GetDefaultImageBytes("default_images/format.png", "png")
	velo2 := &entities.Product{
		Name:        "Велосипед FORMAT 5223 650B (2023, бежевый-мат)",
		Brand:       "FORMAT",
		Description: "Format 5223 27.5 является великолепным вариантом для неспешного катания в городских парках, а также подойдет для ежедневных перемещений по делам. Велосипед оснащён стартовым уровнем оборудования. Рама 27,5' (650B), Tour, 4130 Cr-mo изготовлена из высокотехнологичных и устойчивых к внешним воздействиям материалов. Задний переключатель Microshift Advent M6195M поможет выбрать оптимальный режим для горизонтального движения, а также для подъёма и спуска. Популярная модель велосипедного седла Format All Terrain. Лёгкие покрышки WTB Byway folding, 650x47C с отличным боковым сцеплением позволяют легко маневрировать на дороге. Дисковые механические тормоза Tektro M280, дисковые механические, ротор 160 мм. позволяют мгновенно затормозить перед непредвиденным препятствием. Велосипед оснащён жёсткой вилкой Жесткая, 4130 Cr-mo. Эргономичные и прочные шатуны Prowheel DMA TT, 40T, длина 175 мм. обеспечат эффективное кручение педалей.",
		Image:       velo2_image,
		Price:       1990,
	}
	uc.Products().Create(velo2)

	uc.Products().AssignCategoryByName(2, "Велосипеды")

	uc.Details().CreateValue(1, 2, "Гибридный")
	uc.Details().CreateValue(2, 2, "Мужской")
	uc.Details().CreateValue(3, 2, "Бежевый")
	uc.Details().CreateValue(4, 2, "Хромомолибден")

	// ----------------zam1---------------------------
	zam1_image, _ := GetDefaultImageBytes("default_images/longus.png", "png")
	zam1 := &entities.Product{
		Name:        "Замок на ключ Longus 3550011 6/1800мм Special Edition",
		Brand:       "Longus",
		Description: "Велозамок с креплением на раму или подседельный штырь Ø 27,2 - 31,6 мм. Замок под ключ, 2 ключа в комплекте.",
		Image:       zam1_image,
		Price:       12,
	}
	uc.Products().Create(zam1)

	uc.Products().AssignCategoryByName(3, "Замки велосипедные")

	uc.Details().CreateValue(5, 3, "142г")
	uc.Details().CreateValue(6, 3, "Черный")
	uc.Details().CreateValue(7, 3, "Троссовый")
	uc.Details().CreateValue(8, 3, "На ключ")

	// ----------------zam2---------------------------
	zam2_image, _ := GetDefaultImageBytes("default_images/spencer.png", "png")
	zam2 := &entities.Product{
		Name:        "Замок на ключ U-Lock Spencer ZS2921 165-250мм",
		Brand:       "U-Lock",
		Description: "Замок скоба на ключ для велосипеда, изготовлен из высокопрочной стали, 2 ключа. Крепление для велосипеда в комплекте.",
		Image:       zam2_image,
		Price:       66,
	}
	uc.Products().Create(zam2)

	uc.Products().AssignCategoryByName(4, "Замки велосипедные")

	uc.Details().CreateValue(5, 4, "205г")
	uc.Details().CreateValue(6, 4, "Черный")
	uc.Details().CreateValue(7, 4, "U-образный")
	uc.Details().CreateValue(8, 4, "На ключ")

	// ----------------nas1---------------------------
	nas1_image, _ := GetDefaultImageBytes("default_images/beto.png", "png")
	nas1 := &entities.Product{
		Name:        "Насос BETO CMP-151SG1 (титан)",
		Brand:       "BETO",
		Description: "Насос напольный ручной с манометром.",
		Image:       nas1_image,
		Price:       54,
	}
	uc.Products().Create(nas1)

	uc.Products().AssignCategoryByName(5, "Насосы для велосипеда")

	uc.Details().CreateValue(9, 5, "Насос ручной напольный")
	uc.Details().CreateValue(10, 5, "1035г")
	uc.Details().CreateValue(11, 5, "Черный")

	// ----------------nas2---------------------------
	nas2_image, _ := GetDefaultImageBytes("default_images/moon.png", "png")
	nas2 := &entities.Product{
		Name:        "Насос MOON AIOLOS DIGIT HP/HV 320PSI",
		Brand:       "MOON",
		Description: "Алюминиевый насос высокого давления (до 27 атмосфер) с манометром и шлангом, для подкачки воздушных вилок и задних амортизаторов.",
		Image:       nas2_image,
		Price:       155,
	}
	uc.Products().Create(nas2)

	uc.Products().AssignCategoryByName(6, "Насосы для велосипеда")

	uc.Details().CreateValue(9, 6, "Для амортизаторов")
	uc.Details().CreateValue(10, 6, "560г")
	uc.Details().CreateValue(11, 6, "Графитовый")

	// ----------------sht1---------------------------
	sht1_image, _ := GetDefaultImageBytes("default_images/polisport.png", "png")
	sht1 := &entities.Product{
		Name:        "Щитки для велосипеда Polisport Expander City 8624900001 (чёрный)",
		Brand:       "Polisport",
		Description: "Щитки с дополнительным креплением, предназначены для велосипедов с 26 до 28' колесами. Подходят как для круглых, так и для овальных труб рамы велосипеда. Адаптированы под V-Brake и дисковые тормоза. Удлиняющаяся конструкция крыльев позволяет отрегулировать длину максимально удобным образом, учитывая размер колес и местность для катания.",
		Image:       sht1_image,
		Price:       68,
	}
	uc.Products().Create(sht1)

	uc.Products().AssignCategoryByName(7, "Щитки для велосипеда")

	uc.Details().CreateValue(12, 7, "Комплект")
	uc.Details().CreateValue(13, 7, "537г")
	uc.Details().CreateValue(14, 7, "Черный")
	uc.Details().CreateValue(15, 7, "28")
}
