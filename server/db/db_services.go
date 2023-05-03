package db

import (
	"gorm.io/gorm"
)

var db *gorm.DB

func Transfer(connection *gorm.DB) {
	db = connection
}

func CreateRecord(data interface{}) error {

	err := db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func FindById(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	err := db.Where(column, id).First(data).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateRecord(data interface{}, id interface{}, columName string) *gorm.DB {
	column := columName + "=?"
	result := db.Where(column, id).Updates(data)

	return result
}

func QueryExecutor(query string, data interface{}, args ...interface{}) error {

	err := db.Raw(query, args...).Scan(data).Error
	if err != nil {
		return err
	}

	// return nil if there were no errors
	return nil
}

func DeleteRecord(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	result := db.Where(column, id).Delete(data)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func Delete(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	err := db.Where(column, id).Unscoped().Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func RecordExist(tableName string, columnName string, value string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM " + tableName + " WHERE " + columnName + " = '" + value + "')"
	db.Raw(query).Scan(&exists)
	return exists
}

func BothExists(tablename string, column1 string, value1 string, column2 string, value2 string) bool {
	var exists bool
	query := "select exists(select * from " + tablename + " where " + column1 + " = '" + value1 + "' and " + column2 + " = ' " + value2 + "');"
	db.Raw(query).Scan(&exists)
	return exists
}
