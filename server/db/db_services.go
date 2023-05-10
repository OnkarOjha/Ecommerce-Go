package db

import (
	"fmt"

	"gorm.io/gorm"
)

var db *gorm.DB

func Transfer(connection *gorm.DB) {
	db = connection
}

// Create DB record
func CreateRecord(data interface{}) error {

	err := db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

// Find DB record by ID
func FindById(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	err := db.Where(column, id).Find(data).Error
	if err != nil {
		return err
	}
	return nil
}

//Update Record by ID
func UpdateRecord(data interface{}, id interface{}, columName string) *gorm.DB {
	column := columName + "=?"
	result := db.Where(column, id).Updates(data)

	return result
}

// Execute the provided query
func QueryExecutor(query string, data interface{}, args ...interface{}) error {

	err := db.Raw(query, args...).Scan(data).Error
	if err != nil {
		return err
	}
	return nil
}

//soft delete
func DeleteRecord(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	result := db.Where(column, id).Delete(data)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

// hard delete
func Delete(data interface{}, id interface{}, columName string) error {
	column := columName + "=?"
	err := db.Where(column, id).Unscoped().Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

// Check if given request exists or not
func RecordExist(tableName string, columnName string, value string) bool {
	var exists bool
	query := "SELECT EXISTS(SELECT * FROM " + tableName + " WHERE " + columnName + " = '" + value + "')"
	db.Raw(query).Scan(&exists)
	return exists
}

// check if two matching record exists or not
func BothExists(tablename string, column1 string, value1 string, column2 string, value2 string) bool {
	var exists bool
	query := "select exists(select * from " + tablename + " where " + column1 + " = '" + value1 + "' and " + column2 + " = ' " + value2 + "');"
	db.Raw(query).Scan(&exists)
	return exists
}

//DB function to clear search history table
func SearchHistoryClear() {
	fmt.Println("clearing previous history...")
	db.Exec(`CREATE OR REPLACE FUNCTION truncate_rows_if_exceeds_threshold(threshold INTEGER, table_name TEXT) RETURNS VOID AS $$
	DECLARE
		count_rows INTEGER;
	BEGIN
		-- Get the count of rows in the table
		EXECUTE format('SELECT COUNT(*) FROM %I', table_name) INTO count_rows;
	
		-- Check if the count of rows exceeds the threshold value
		IF count_rows > threshold THEN
			-- Truncate the first 10 rows
			
			EXECUTE format('DELETE FROM %I WHERE id IN (SELECT id FROM %I LIMIT 10)', table_name, table_name);
		END IF;
	END;
	$$ LANGUAGE plpgsql;
	`)
	db.Exec(`SELECT truncate_rows_if_exceeds_threshold(20, 'search_histories');`)
}
