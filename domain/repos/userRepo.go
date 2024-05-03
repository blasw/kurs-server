package repos

import (
	"errors"
	"gorm.io/gorm"
	"kurs-server/domain/entities"
)

type UserRepo struct {
	Storage *gorm.DB
}

// Create TODO: Might work incorrect
func (u *UserRepo) Create(newUser *entities.User) (uint, error) {
	tx := u.Storage.Create(newUser)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return newUser.ID, nil
}

func (u *UserRepo) Delete(user *entities.User) error {
	return nil
}

func (u *UserRepo) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	tx := u.Storage.Where("username = ?", username).First(&user)
	return &user, tx.Error
}

func (u *UserRepo) GetUserByRefreshToken(refreshToken string) (*entities.User, error) {
	if refreshToken == "" {
		return nil, errors.New("invalid token")
	}

	var user *entities.User
	tx := u.Storage.First(&user, "refresh_token", refreshToken)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (u *UserRepo) UpdateUserRefreshToken(username string, newRefreshToken string) error {
	var user *entities.User
	tx := u.Storage.First(&user, "username", username)
	if tx.Error != nil {
		return tx.Error
	}

	tx = u.Storage.Model(&user).Update("refresh_token", newRefreshToken)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
