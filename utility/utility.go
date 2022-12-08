package utility

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

func ContextJSON(context *gin.Context, httpStatus int, result interface{}, error error) {
	context.JSON(httpStatus, gin.H{
		"error":  error,
		"result": result,
	})
}

func RequestBody(context *gin.Context) ([]byte, error) {
	return io.ReadAll(context.Request.Body)
}

func RawBodyToMap(rawJSON []byte) map[string]interface{} {
	dictionary := map[string]interface{}{}
	json.Unmarshal(rawJSON, &dictionary)
	return dictionary
}

func RequestBodyToMap(context *gin.Context) map[string]interface{} {
	rawJSON, _ := RequestBody(context)
	return RawBodyToMap(rawJSON)
}
