package main

import (
	"STShortURL/database"
	"STShortURL/log"
	"STShortURL/utility"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	db := database.Connect("./STShortURL.sqlite3")
	if db != nil {
		router := gin.Default()
		insertShortUrl(router, db)
		selectShortUrl(router, db)
		err := router.Run(":8080")
		if err != nil {
			log.Error.Println(err)
		}
	}
}

func insertShortUrl(router *gin.Engine, db *gorm.DB) {
	router.POST("/", func(context *gin.Context) {
		dictionary := utility.RequestBodyToMap(context)
		result, err := database.Insert(db, dictionary)
		utility.ContextJSON(context, http.StatusOK, result, err)
	})
}

func selectShortUrl(router *gin.Engine, db *gorm.DB) {
	router.GET("/:code", func(context *gin.Context) {
		notfoundPage := "https://smallten.me"
		code := context.Param("code")
		fmt.Println("Code", code)
		result := database.Select(db, code)
		refreshHtml := fmt.Sprintf("<html><meta http-equiv='refresh' content='0;url=%s'/></html>", notfoundPage)
		defer func() {
			context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(refreshHtml))
		}()
		if result.ID == 0 {
			return
		}
		refreshHtml = fmt.Sprintf("<html><meta http-equiv='refresh' content='0;url=%s'/></html>", result.Url)
	})
}
