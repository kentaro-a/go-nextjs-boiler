package middleware

import (
	"gorm.io/gorm"
)

type Middleware struct {
	DB *gorm.DB
}
