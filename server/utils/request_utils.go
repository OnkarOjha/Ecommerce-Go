package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

// Decoding the request body
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

// Set Header for content
func SetHeader(context *gin.Context) {
	context.Header("Access-Control-Allow-Origin", "*")
	context.Header("Content-Type", "application/json")

}
