package repository

import (
	"GoGin-API-Base/dao"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByEmail(email string) (dao.User, error)
	FindUserById(id int) (dao.User, error)
	Save(user *dao.User) (dao.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (u UserRepositoryImpl) FindUserByEmail(email string) (dao.User, error) {
	var user dao.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Error("User not found. Error: ", err)
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) FindUserById(id int) (dao.User, error) {
	user := dao.User{
		ID: id,
	}
	err := u.db.First(&user).Error
	if err != nil {
		log.Error("Got and error when find user by id. Error: ", err)
		return dao.User{}, err
	}
	return user, nil
}

func (u UserRepositoryImpl) Save(user *dao.User) (dao.User, error) {
	err := u.db.Create(&user).Error
	if err != nil {
		log.Error("User not created. Error: ", err)
		processedError := ProcessError(err)
		return dao.User{}, processedError
	}
	return *user, nil
}

func ProcessError(err error) error {
	pgErrCode := err.(*pgconn.PgError).Code
	processedError := err

	if pgErrCode == "23505" {
		processedError = errors.New("the email or the user is already in use")
	}

	return processedError
}

func UserRepositoryInit(db *gorm.DB) *UserRepositoryImpl {
	db.AutoMigrate(&dao.User{})
	return &UserRepositoryImpl{
		db: db,
	}
}
