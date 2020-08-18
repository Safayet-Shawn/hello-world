package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)
type Register struct{
	ID       int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Email string `gorm:"type:Varchar(100)" json:"Email,omitempty"`
	Name  string `gorm:"type:varchar(50)"json:"Name,omitempty"`
	Phone string `gorm:"varchar(20);unique"json:"Phone,omitempty"`
	Password  string `gorm:"varchar(100)"json:"Password,omitempty"`
}
type Login struct{
	ID       int64 `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Email   string `gorm:"type:Varchar(100)" json:"Email,omitempty"`
	Password   string `gorm:"varchar(100)"json:"Password,omitempty"`
}
 var db *gorm.DB
func initDb(){
	var err error
	db,err:= gorm.Open("mysql","root:itsshawn@007@@tcp(localhost:3306)/?parseTime=True")
	if err !=nil{
		fmt.Println(err)
		panic("failed to connect Database")
	}
	db.Exec("CREATE DATABASE loginregs")
	
	db.Exec("use loginregs")
	db.AutoMigrate(&Login{},&Register{})
}
func regUser(c echo.Context)error{
	var reg Register
	defer c.Request().Body.Close()
	err:=json.NewDecoder(c.Request().Body).Decode(&reg)
	if err !=nil{
		log.Printf("failed processing request: %s",err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Register : %#v",reg)
	return c.JSON(http.StatusCreated, reg)

	// return c.String(http.StatusOK,"You are registered")
}
func main() {
		initDb()
	e := echo.New()

	//router
	e.POST("/register",regUser)
	// e.POST("/login_user", loginUser)
	// e.GET("/user/:id", setUser)
	// e.GET("/users", all_User) 
	e.Logger.Fatal(e.Start(":8080"))
}

