package model

import (
	"time"

	_ "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MailAuthModel struct {
	DB *gorm.DB
}

type MailAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id" ja:"ID"`
	Function  string    `json:"function" validate:"required,max=255" ja:"機能"`
	Mail      string    `gorm:"unique" json:"mail" validate:"email,required,max=255" ja:"メールアドレス"`
	Token     string    `json:"token" validate:"required,max=255" ja:"トークン"`
	StatusFlg int       `json:"status_flg" validate:"number,min=0,max=9" ja:"ステータス"`
	ExpireAt  time.Time `gorm:"expire_at" json:"expire_at"`
	CreatedAt time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at" json:"updated_at"`
}

func NewMailAuthModel(db *gorm.DB) *MailAuthModel {
	return &MailAuthModel{
		DB: db,
	}
}

func (m *MailAuthModel) Create(user_mail_auth *MailAuth) error {
	err := m.DB.Create(user_mail_auth).Error
	return errors.WithStack(err)
}

func (m *MailAuthModel) DeleteByMailFunction(mail string, function string) error {
	err := m.DB.Where("mail = ? AND function = ?", mail, function).Delete(&MailAuth{}).Error
	return errors.WithStack(err)
}

func (m *MailAuthModel) FindByMailFunction(mail string, function string) ([]MailAuth, error) {
	list := []MailAuth{}
	err := m.DB.Where("mail = ? AND function = ?", mail, function).Find(&list).Error
	return list, errors.WithStack(err)
}

func (m *MailAuthModel) IsExpiredToken(user_mail_auth *MailAuth) (bool, error) {
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
