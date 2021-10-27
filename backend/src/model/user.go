package model

import (
	"app/config"
	"crypto/sha256"
	"encoding/hex"
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserModel struct {
	DB *gorm.DB
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id" ja:"ID"`
	Name      string    `json:"name" validate:"required,max=255" ja:"ユーザー名"`
	Mail      string    `gorm:"unique" json:"mail" validate:"email,required,max=255" ja:"メールアドレス"`
	Password  string    `json:"password" validate:"required,min=8,max=64" ja:"パスワード"`
	StatusFlg int       `json:"status_flg" validate:"number,min=0,max=9" ja:"ステータス"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}

func (user User) NewUserWithoutSecrets() User {
	user.Password = ""
	return user
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		DB: db,
	}
}

func (m *UserModel) GetHashedPassword(password string) string {
	c := config.Get()
	s := c.Salt + password
	r := sha256.Sum256([]byte(s))
	return hex.EncodeToString(r[:])
}

func (m *UserModel) Create(user *User) error {
	err := m.DB.Create(user).Error
	return errors.WithStack(err)
}

func (m *UserModel) Save(user *User) error {
	err := m.DB.Save(user).Error
	return errors.WithStack(err)
}

func (m *UserModel) IsMailExist(mail string) (bool, *User, error) {
	is_exist := false
	user := User{}
	var err error

	err = m.DB.Where(
		"mail = ? AND status_flg = ?",
		mail,
		0,
	).First(&user).Error

	if err == nil {
		is_exist = true
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// not found は正常
			err = nil
		}
	}
	return is_exist, &user, errors.WithStack(err)
}

func (m *UserModel) FindByMail(mail string) (*User, error) {
	var user User
	// status_flg = 0などstructで検索するとゼロ値は検索条件から除外されるので直接指定する
	err := m.DB.Where(
		"mail = ? AND status_flg = ?",
		mail,
		0,
	).First(&user).Error
	return &user, errors.WithStack(err)
}

func (m *UserModel) FindByMailPassword(mail, encrypted_password string) (*User, error) {
	var user User
	// status_flg = 0などstructで検索するとゼロ値は検索条件から除外されるので直接指定する
	err := m.DB.Where(
		"mail = ? AND password = ? AND status_flg = ?",
		mail,
		encrypted_password,
		0,
	).First(&user).Error
	return &user, errors.WithStack(err)
}
