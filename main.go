package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "time"//
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

//Register struct//
type Register struct {
	// gorm.Model
	ID       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT;"`
	Email    string `json:"email,omitempty" gorm:"type:varchar(100);" `
	Name     string `json:"name,omitempty" gorm:"type:varchar(50);" `
	Phone    string `json:"phone,omitempty" gorm:"type:varchar(20);unique;" `
	Password string `json:"password,omitempty" gorm:"type:varchar(100);" `
}

//Login struct//
type Login struct {
	// gorm.Model
	 LID       int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT;"`
	Email    string `json:"Email,omitempty" gorm:"type:varchar(100);" `
	Password string `json:"Password,omitempty" gorm:"type:varchar(100);"`
}

var db *gorm.DB

func initDb() {
	var err error
	db, err = gorm.Open("mysql", "root:itsshawn@007@@tcp(localhost:3306)/logrreg?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect Database")
	}
	 // db.Exec("CREATE DATABASE signn")

	 db.Exec("use logrreg")
	db.AutoMigrate(&Login{}, &Register{})
}
func regUser(c echo.Context) error {
	reg :=new(Register)
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reg)
	if err != nil {
		log.Printf("failed processing request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Register : %#v", reg)
	err=db.Table("registers").Create(reg).Error
	if err!= nil{
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, reg)
	// return c.String(http.StatusOK,"You are registered")
}
func main() {
	initDb()
	e := echo.New()

	//router
	e.POST("/register", regUser)
	// e.POST("/login_user", loginUser)
	// e.GET("/user/:id", setUser)
	// e.GET("/users", all_User)
	e.Logger.Fatal(e.Start(":8080"))
}
