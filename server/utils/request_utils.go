package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func RequestDecoding(context *gin.Context, data interface{}) error {

	reqBody, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(reqBody, &data)
	if err != nil {
		return err
	}
	return nil
}

func SetHeader(context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/json")

}
