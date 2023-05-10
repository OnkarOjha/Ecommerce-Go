package db

import (
	"main/server/model"

	"gorm.io/gorm"
)

// Auto Migrate DB
func AutoMigrateDatabase(db *gorm.DB) {

	var dbVersion model.DbVersion
	err := db.First(&dbVersion).Error
	if err != nil {
		panic(err)
	}
	if dbVersion.Version < 1 {
		err := db.AutoMigrate(&model.User{}, &model.Session{}, &model.Cart{}, &model.CartProducts{}, &model.Order{}, &model.UserAddresses{}, &model.UserPayments{}, &model.Payment{}, &model.SearchHistory{}, &model.DbConstant{}, &model.Vendor{}, &model.Products{}, &model.VendorInventory{}, &model.Coupons{}, &model.CouponRedemptions{}, &model.OrderRequest{})
		if err != nil {
			panic(err)
		}
		db.Create(&model.DbVersion{
			Version: 1,
		})
		dbVersion.Version = 1
	}

}
