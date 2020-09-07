package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

//Register struct//
type Register struct {
	// gorm.Model
	ID    int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;"`
	Email string `json:"email,omitempty" gorm:"type:varchar(100);unique;" `
	Name  string `json:"name,omitempty" gorm:"type:varchar(50);" `
	Phone string `json:"phone,omitempty" gorm:"type:varchar(20);" `
	// Password     string `json:"password,omitempty" gorm:"type:varchar(100);"`
	PasswordHash string `json:"passwordhash,omitempty" gorm:"type:varchar(100);" `
}

var db *gorm.DB

func initDb() {
	var err error
	db, err = gorm.Open("mysql", "root:itsshawn@007@@tcp(localhost:3306)/shawn?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect Database")
	}
	//db.Exec("CREATE DATABASE signin")

	db.Exec("use shawn")
	db.AutoMigrate(&Register{})
}

//User function
func User(c echo.Context) error {
	var user []Register

	err := db.Find(&user)
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, user)
}
func userByID(c echo.Context) error {
	userr := new(Register)
	ID := c.QueryParam("id")
	err := db.Where("id = ? ", ID).First(&userr).Error

	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusOK, userr)

}
func main() {
	initDb()
	e := echo.New()

	//router

	e.GET("api/v1/auth/users", User)
	e.GET("api/v1/auth/userid", userByID)

	e.Logger.Fatal(e.Start(":8081"))

}
