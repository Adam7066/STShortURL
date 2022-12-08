package main

import (
	"STShortURL/database"
	"STShortURL/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	db := database.Connect("./DB.sqlite3")
	if db != nil {
		router := gin.Default()
		insertShortUrl(router, db)
		selectShortUrl(router, db)
		router.Run(":8080")
	}
}

func insertShortUrl(router *gin.Engine, db *gorm.DB) {
	router.POST("/shortURL", func(context *gin.Context) {
		dictionary := utility.RequestBodyToMap(context)
		result, error := database.Insert(db, dictionary)
		utility.ContextJSON(context, http.StatusOK, result, error)
	})
}

func selectShortUrl(router *gin.Engine, db *gorm.DB) {
	router.GET("/shortURL/:code", func(context *gin.Context) {
		blog := "https://blog.smallten.me"
		code := context.Param("code")
		result := database.Select(db, code)
		refreshHtml := fmt.Sprintf("<html><meta http-equiv='refresh' content='0;url=%s'/></html>", blog)
		defer func() {
			context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(refreshHtml))
		}()
		if result.ID == 0 {
			return
		}
		refreshHtml = fmt.Sprintf("<html><meta http-equiv='refresh' content='0;url=%s'/></html>", result.Url)
	})
}
