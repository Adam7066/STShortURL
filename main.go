package main

import (
	"STShortURL/database"
	"STShortURL/log"
	"STShortURL/utility"
	"fmt"
	"github.com/pquerna/otp/totp"
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
		accountName := dictionary["account_code"].(string)
		passcode := dictionary["passcode"].(string)
		if validateTOTP(db, accountName, passcode) {
			url := dictionary["url"].(string)
			result, err := database.Insert(db, url)
			utility.ContextJSON(context, http.StatusOK, result, err)
		} else {
			result := map[string]interface{}{"isSuccess": false}
			utility.ContextJSON(context, http.StatusUnauthorized, result, nil)
		}
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

func validateTOTP(db *gorm.DB, accountName, passcode string) bool {
	secret := database.SelectTotp(db, accountName)
	return totp.Validate(passcode, secret)
}
