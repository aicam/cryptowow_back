package Bridge

import "gorm.io/gorm"

type Server struct {
	DB            *gorm.DB
	BetOperations struct {
	}
}
