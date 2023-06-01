package db

import (
	"main/server/model"

	"gorm.io/gorm"
)

// DB executables function
func Execute(db *gorm.DB) {
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	err := db.AutoMigrate(&model.DbVersion{})
	if err != nil {
		return
	}
}
