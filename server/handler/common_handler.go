package handler

import (
	"main/server/services/common_services"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary  		User Register Handler
// @Description  	Registering User with initial details in DB
// @Tags 			User
// @Accept 			json
// @Procedure 		json
// @Param   		username query string true "Username" SchemaExample({ "username" : "string"})
// @Success			200	{string}	response.Response
// @Failure			400	{string}	response.Response
// @Failure			409	{string}	response.Response
// @Failure			500	{string}	response.Response
// @Router			/username-exists [get]
func UserNameHandler(ctx *gin.Context) {
	utils.SetHeader(ctx)

	common_services.UserNameCheck(ctx)
}
