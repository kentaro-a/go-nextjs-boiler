package health_check

import (
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}
