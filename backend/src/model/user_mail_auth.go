package model

import (
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserMailAuthModel struct {
	DB *gorm.DB
}

type UserMailAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id" ja:"ID"`
	Function  string    `json:"function" validate:"required,max=255" ja:"機能"`
	Mail      string    `gorm:"unique" json:"mail" validate:"email,required,max=255" ja:"メールアドレス"`
	Token     string    `json:"token" validate:"required,max=255" ja:"トークン"`
	StatusFlg int       `json:"status_flg" validate:"number,min=0,max=9" ja:"ステータス"`
	ExpireAt  time.Time `gorm:"expire_at" json:"expire_at"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}

func NewUserMailAuthModel(db *gorm.DB) *UserMailAuthModel {
	return &UserMailAuthModel{
		DB: db,
	}
}

func (m *UserMailAuthModel) Create(user_mail_auth *UserMailAuth) error {
	err := m.DB.Create(user_mail_auth).Error
	return errors.WithStack(err)
}

func (m *UserMailAuthModel) DeleteByMailFunction(mail string, function string) error {
	err := m.DB.Where("mail = ? AND function = ?", mail, function).Delete(&UserMailAuth{}).Error
	return errors.WithStack(err)
}

func (m *UserMailAuthModel) FindByMailFunction(mail string, function string) ([]UserMailAuth, error) {
	list := []UserMailAuth{}
	err := m.DB.Where("mail = ? AND function = ?", mail, function).Find(&list).Error
	return list, errors.WithStack(err)
}

func (m *UserMailAuthModel) IsExpiredToken(user_mail_auth *UserMailAuth) (bool, error) {
	err := m.DB.Where("token = ? AND status_flg = 0 AND expire_at > ? ", user_mail_auth.Token, time.Now()).First(user_mail_auth).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, errors.WithStack(err)
		} else {
			return true, nil
		}
	}
	return false, nil
}
