package filter

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func FilterByCategoryService(context *gin.Context) {
	category := context.Query("category")

	if category == "" {
		response.ErrorResponse(context, 400, "No product Category specified")
		return
	}

	if !db.RecordExist("products", "product_category", category) {
		response.ErrorResponse(context, 400, "Product Category does not exist")
		return
	}

	var productByCategory []model.Products

	query := "SELECT * FROM products where product_category='" + strings.ToLower(category) + "' ORDER BY product_price DESC LIMIT 30"
	err := db.QueryExecutor(query, &productByCategory)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Finding in DB")
		return
	}

	response.ShowResponse(
		"Success",
		200,
		"Here are the list of products according to the given category",
		productByCategory,
		context,
	)
}

func FilterByPriceService(context *gin.Context) {
	priceFrom := context.Query("from")
	priceTo := context.Query("to")

	if priceFrom == "" || priceTo == "" {
		response.ErrorResponse(context, 400, "No product Price specified")
		return
	}

	priceFromInt, _ := strconv.Atoi(priceFrom)
	priceToInt, _ := strconv.Atoi(priceTo)

	if priceFromInt < 0 {
		response.ErrorResponse(context, 400, "Starting Price must be greater than zero")
		return
	}
	if priceToInt < 0 {
		response.ErrorResponse(context, 400, "Starting Price must be greater than zero")
		return
	}

	var productByPrice []model.Products

	query := "select * from products where product_price BETWEEN " + priceFrom + " AND " + priceTo + " ORDER BY product_price ASC LIMIT 30;"
	err := db.QueryExecutor(query, &productByPrice)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Finding in DB")
		return
	}

	response.ShowResponse(
		"Success",
		200,
		"Here are the list of products according to the given price range",
		productByPrice,
		context,
	)

}

func FilterByBrandService(context *gin.Context) {
	brandName := context.Query("brand")

	if brandName == "" {
		response.ErrorResponse(context, 400, "No brand name specified")
		return
	}

	if !db.RecordExist("products", "product_brand", brandName) {
		response.ErrorResponse(context, 400, "Product brand does not exist")
		return
	}

	var productByBrand []model.Products

	query := "SELECT * FROM products where product_category='" + strings.ToLower(brandName) + "' ORDER BY product_price DESC LIMIT 30"
	err := db.QueryExecutor(query, &productByBrand)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Finding in DB")
		return
	}

	response.ShowResponse(
		"Success",
		200,
		"Here are the list of products according to the given category",
		productByBrand,
		context,
	)

}

func SearchBarService(context *gin.Context) {
	productQuery := context.Query("productQuery")

	var productNameSearch []model.Products

	var productNameSearchExists bool

	boolQuery := "SELECT EXISTS (select * from products where product_name LIKE LOWER('%\\" + productQuery + "%'))"

	err := db.QueryExecutor(boolQuery, &productNameSearchExists)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Finding in DB")
		return
	}

	if productNameSearchExists {
		productNameSearchQuery := "select * from products where product_name LIKE LOWER('%\\" + productQuery + "%')  ORDER BY product_price ASC LIMIT 10"

		err := db.QueryExecutor(productNameSearchQuery, &productNameSearch)
		if err != nil {
			response.ErrorResponse(context, 500, "Error Finding in DB")
			return
		}

		// search with price range
		priceFrom := context.Query("from")
		priceTo := context.Query("to")
		currentTime := time.Now()

		if priceFrom != "" && priceTo != "" {
			SearchWithPriceRange(context, productNameSearch, priceFrom, priceTo)
			SearchHistoryUpdate(context, productNameSearch, productQuery, priceFrom, priceTo, currentTime)
			return
		}

		if priceFrom != "" {
			SearchWithPriceFromRange(context, productNameSearch, priceFrom)
			SearchHistoryUpdate(context, productNameSearch, productQuery, priceFrom, priceTo, currentTime)
			return
		}

		if priceTo != "" {
			SearchWithPriceToRange(context, productNameSearch, priceTo)
			SearchHistoryUpdate(context, productNameSearch, productQuery, priceFrom, priceTo, currentTime)
			return
		}

		response.ShowResponse(
			"Success",
			200,
			"The List of products are",
			productNameSearch,
			context,
		)

		SearchHistoryUpdate(context, productNameSearch, productQuery, priceFrom, priceTo, currentTime)

	} else {
		response.ErrorResponse(context, 400, "Product with this name doesn't exist")
		return
	}

}

func SearchWithPriceRange(context *gin.Context, productNameSearch []model.Products, priceFrom string, priceTo string) {

	priceFromInt, _ := strconv.Atoi(priceFrom)
	priceToInt, _ := strconv.Atoi(priceTo)

	if priceFromInt < 0 {
		response.ErrorResponse(context, 400, "Starting Price must be greater than zero")
		return
	}
	if priceToInt < 0 {
		response.ErrorResponse(context, 400, "Starting Price must be greater than zero")
		return
	}

	var productByPrice []model.Products

	for i := range productNameSearch {
		if productNameSearch[i].ProductPrice >= float64(priceFromInt) && productNameSearch[i].ProductPrice <= float64(priceToInt) {
			productByPrice = append(productByPrice, productNameSearch[i])
		}
	}
	if productByPrice == nil {
		response.ShowResponse(
			"Success",
			200,
			"The List of products are",
			productNameSearch,
			context,
		)
	}
	response.ShowResponse(
		"Success",
		200,
		"Here are the list of products according to the given price range",
		productByPrice,
		context,
	)

}

func SearchWithPriceFromRange(context *gin.Context, productNameSearch []model.Products, priceFrom string) {
	priceFromInt, _ := strconv.Atoi(priceFrom)

	var productByPrice []model.Products

	for i := range productNameSearch {
		if productNameSearch[i].ProductPrice >= float64(priceFromInt) {
			productByPrice = append(productByPrice, productNameSearch[i])
		}
	}
	if productByPrice == nil {
		response.ShowResponse(
			"Success",
			200,
			"The List of products are",
			productNameSearch,
			context,
		)
	}
	response.ShowResponse(
		"Success",
		200,
		"List of products according from the price specified",
		productByPrice,
		context,
	)
}

func SearchWithPriceToRange(context *gin.Context, productNameSearch []model.Products, priceTo string) {
	priceToInt, _ := strconv.Atoi(priceTo)

	var productByPrice []model.Products

	for i := range productNameSearch {
		if productNameSearch[i].ProductPrice <= float64(priceToInt) {
			productByPrice = append(productByPrice, productNameSearch[i])
		}
	}
	if productByPrice == nil {
		response.ShowResponse(
			"Success",
			200,
			"The List of products are",
			productNameSearch,
			context,
		)
	}
	response.ShowResponse(
		"Success",
		200,
		"List of products according under the price specified",
		productByPrice,
		context,
	)
}

var frequency int

func SearchHistoryUpdate(context *gin.Context, productNameSearch []model.Products, productQuery string, priceFrom string, priceTo string, currentTime time.Time) {
	var searchHistoryUpdate []model.SearchHistory

	frequency++
	var productid string
	for _, p := range productNameSearch {
		productid = p.ProductId
		var temp model.SearchHistory
		temp.ProductId = productid
		temp.SearchTime = currentTime
		temp.SearchFrequency = frequency
		temp.SearchQuery = productQuery + " in between " + priceFrom + " and " + priceTo
		searchHistoryUpdate = append(searchHistoryUpdate, temp)

	}

	// //search history clear before next update
	db.SearchHistoryClear()
	err := db.CreateRecord(&searchHistoryUpdate)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Creating Record")
		return
	}
}

func SearchBarHistoryService(context *gin.Context) {
	var searchBarHistoryLoader []model.SearchHistory
	query := "SELECT DISTINCT search_frequency, search_time ,search_query FROM search_histories ORDER BY search_time ASC LIMIT 5;"

	err := db.QueryExecutor(query, &searchBarHistoryLoader)
	if err != nil {
		response.ErrorResponse(context, 500, "Error Finding in DB")
		return
	}

	var searchQueryResponse []response.SearchResponse
	for _, searchBarHistory := range searchBarHistoryLoader {
		var temp response.SearchResponse
		temp.SearchQuery = searchBarHistory.SearchQuery
		searchQueryResponse = append(searchQueryResponse, temp)
	}

	response.ShowResponse(
		"Success",
		200,
		"Here are the list of search history",
		searchQueryResponse,
		context,
	)
}
