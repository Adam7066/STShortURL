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

type Totp struct {
	gorm.Model
	AccountName string `json:"accountName" gorm:"unique"`
	Secret      string `json:"secret" gorm:"unique"`
}

func Connect(DBPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		log.Error.Println(err)
		return nil
	}
	db.AutoMigrate(&ShortURL{})
	db.AutoMigrate(&Totp{})
	return db
}

func Insert(db *gorm.DB, url string) (map[string]interface{}, error) {
	isSuccess := false
	var shortURLs []ShortURL
	err := db.Where("url=?", url).Limit(1).Find(&shortURLs).Error
	if err != nil {
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

func InsertTotp(db *gorm.DB, accountName string, secret string) error {
	var totps []Totp
	err := db.Where("account_name=?", accountName).Limit(1).Find(&totps).Error
	if err != nil {
		return err
	}
	err = db.Create(&Totp{AccountName: accountName, Secret: secret}).Error
	return err
}

func SelectTotp(db *gorm.DB, accountName string) string {
	var tmpTotp Totp
	db.Take(&tmpTotp, "account_name=?", accountName)
	return tmpTotp.Secret
}
