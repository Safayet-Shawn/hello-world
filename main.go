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
	ID       int64 `json:"Id"`
	Email string `json:"Email"`
	Name  string `json:"Name"`
	Phone string `json:"Phone"`
	Password  string `json:"Password"`
}
type Login struct{
	ID       int64 `json:"Id"`
	Email   string `json:"email"`
	Password   string `json:"password"`
}
 var db *gorm.DB
func initDb(){
	var err error
	db,err:= gorm.Open("mysql","root:itsshawn@007@@tcp(localhost:3306)/?parseTime=True")
	if err !=nil{
		fmt.Println(err)
		panic("failed to connect Database")
	}
	db.Exec("CREATE DATABASE loginregst")
	
	db.Exec("use loginregst")
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

