package handler

import (
	"main/server/services/filter"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary  		Filter By Category
// @Tags 			Search&Filter
// @Accept 			json
// @Procedure 		json
// @Param			category query string true "Category" SchemaExample({"category" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/filter/category [get]
func FilterByCategoryHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByCategoryService(ctx)

}

// @Summary  		Filter By Price
// @Tags 			Search&Filter
// @Accept 			json
// @Procedure 		json
// @Param			from query string true "Price" SchemaExample({"from" : "string"})
// @Param			to query string true "Price" SchemaExample({"to" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/filter/price [get]
func FilterByPriceHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByPriceService(ctx)
}

// @Summary  		Filter By Brand
// @Tags 			Search&Filter
// @Accept 			json
// @Procedure 		json
// @Param			brand query string true "Brand" SchemaExample({"brand" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/filter/brand [get]
func FilterByBrandHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByBrandService(ctx)
}

// @Summary  		Search Bar
// @Tags 			Search&Filter
// @Accept 			json
// @Procedure 		json
// @Param			productQuery query string true "Search Bar" SchemaExample({"productQuery" : "string"})
// @Param			from query string true "Search Bar" SchemaExample({"from" : "string"})
// @Param			to query string true "Search Bar" SchemaExample({"to" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/search-bar [get]
func SearchBarHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.SearchBarService(ctx)
}

// @Summary  		Search Bar History
// @Tags 			Search&Filter
// @Accept 			json
// @Procedure 		json
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/search-bar/history [get]
func SearchBarHistoryHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.SearchBarHistoryService(ctx)
}
