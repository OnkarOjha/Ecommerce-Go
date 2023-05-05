package handler

import (
	"main/server/services/filter"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func FilterByCategoryHandler(context *gin.Context) {
	utils.SetHeader(context)

	filter.FilterByCategoryService(context)

}

func FilterByPriceHandler(context *gin.Context) {
	utils.SetHeader(context)

	filter.FilterByPriceService(context)
}

func FilterByBrandHandler(context *gin.Context) {
	utils.SetHeader(context)

	filter.FilterByBrandService(context)
}

func SearchBarHandler(context *gin.Context) {
	utils.SetHeader(context)

	filter.SearchBarService(context)
}

func SearchBarHistoryHandler(context *gin.Context) {
	utils.SetHeader(context)

	filter.SearchBarHistoryService(context)
}
