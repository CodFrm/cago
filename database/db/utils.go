package db

import "gorm.io/gorm"

func RecordNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
