package database

import (
	"STShortURL/log"
	"crypto/md5"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ShortURL struct {
	gorm.Model
	Url  string `json:"url" gorm:"index:id, unique"`
	Code string `json:"code" gorm:"index:id, unique"`
}

func Connect(DBPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		log.Error.Println(err)
		return nil
	}
	db.AutoMigrate(&ShortURL{})
	return db
}

func Insert(db *gorm.DB, dictionary map[string]interface{}) (map[string]interface{}, error) {
	isSuccess := false
	url := dictionary["url"].(string)
	var shortURLs []ShortURL
	err := db.Where("url=?", url).Limit(1).Find(&shortURLs).Error
	if err == nil {
		return nil, err
	}
	code := fmt.Sprintf("%x", md5.Sum([]byte(url)))[:6]
	err = db.Create(&ShortURL{Url: url, Code: code}).Error
	if err == nil {
		isSuccess = true
	}
	result := map[string]interface{}{"isSuccess": isSuccess, "code": code}
	return result, err
}

func Select(db *gorm.DB, code string) ShortURL {
	var tmpShortURL ShortURL
	db.Take(&tmpShortURL, "code=?", code)
	return tmpShortURL
}
