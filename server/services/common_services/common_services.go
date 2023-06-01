package common_services

import (
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

func UserNameCheck(ctx *gin.Context) {
	var usersnameexists bool
	var users []model.User
	username := ctx.Query("username")
	query := "SELECT EXISTS(select * from users where user_name LIKE '%" + username + "%');"
	err := db.QueryExecutor(query, &usersnameexists)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Cannot execute query")
		return
	}
	var matchUserName []string
	if usersnameexists {

		query2 := "select * from users where user_name LIKE '%" + username + "%';"
		err := db.QueryExecutor(query2, &users)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_INTERNAL_SERVER_ERROR, "Unable to execute query")
			return
		}
		for _, user := range users {
			matchUserName = append(matchUserName, user.UserName)
		}
		response.ShowResponse(
			"Success",
			utils.HTTP_OK,
			"User with this substring exists",
			matchUserName,
			ctx,
		)
		return
	} else {
		response.ShowResponse("Success",
			utils.HTTP_OK,
			"",
			matchUserName,
			ctx,
		)
		return
	}
}
