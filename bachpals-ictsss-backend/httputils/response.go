package httputils

import "github.com/gin-gonic/gin"

type ResponseType struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func AbortWithStatusJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	responseStruct := ResponseType{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	c.AbortWithStatusJSON(responseStruct.StatusCode, responseStruct)
	return
}

func ResponseJson(c *gin.Context, statusCode int, message string, data interface{}) {
	responseStruct := ResponseType{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	c.JSON(responseStruct.StatusCode, responseStruct)
	return
}
