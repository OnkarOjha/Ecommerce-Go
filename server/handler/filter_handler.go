package handler

import (
	"main/server/services/filter"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func FilterByCategoryHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByCategoryService(ctx)

}

func FilterByPriceHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByPriceService(ctx)
}

func FilterByBrandHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.FilterByBrandService(ctx)
}

func SearchBarHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.SearchBarService(ctx)
}

func SearchBarHistoryHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	filter.SearchBarHistoryService(ctx)
}
