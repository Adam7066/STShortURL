package main

import (
	"STShortURL/database"
	"STShortURL/utility"
	"flag"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
	"image/png"
	"os"
)

func main() {
	var accountName string
	flag.StringVar(&accountName, "ac", "", "accountName")
	flag.Parse()
	if accountName != "" {
		if db := database.Connect("./STShortURL.sqlite3"); db != nil {
			generateTOTP(db, accountName)
		}
	}
}

func generateTOTP(db *gorm.DB, accountName string) {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "STShortURL",
		AccountName: accountName,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	// save secret key to database
	utility.CheckError(database.InsertTotp(db, accountName, key.Secret()))
	// save QR code
	f, _ := os.Create("totp.png")
	defer f.Close()
	img, err := key.Image(200, 200)
	if err != nil {
		panic(err)
	}
	utility.CheckError(png.Encode(f, img))
}
