package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/middleware"
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



var db *gorm.DB

func initDb() {
	var err error
	db, err = gorm.Open("mysql", "root:itsshawn@007@@tcp(localhost:3306)/sign?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect Database")
	}
	 // db.Exec("CREATE DATABASE signn")

	 db.Exec("use sign")
	db.AutoMigrate( &Register{})
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
func loginUser(c echo.Context)error{
	lgp :=new (Register)
	Email:=c.QueryParam("email")
	Password:=c.QueryParam("password")
	err:=db.Where("email = ? AND password= ? ", Email,Password).First(&lgp).Error
	if err!=nil{
		log.Println(err)
	} else{
	//create tooken
	tk:=jwt.New(jwt.SigningMethodHS256)
	claims:=tk.Claims.(jwt.MapClaims)
	claims["name"]=lgp.Name
	claims["email"]=lgp.Email
	claims["phone"]=lgp.Phone
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	
	tkn, err := tk.SignedString([]byte("secret"))
        if err != nil {
            return err
        }
        return c.JSON(http.StatusOK, map[string]string{

            "token": tkn,
        })
	// fmt.Println(lgp.Email,lgp.Password)
	// return c.JSON(http.StatusCreated, lgp)
}
	return echo.ErrUnauthorized
}
func jwtLogin(c echo.Context)error{
	return c.String(http.StatusOK,"You are Logged in sucessfully !")
}
func main() {
	initDb()
	e := echo.New()

	//jwt group//
	jwtGroup:=e.Group("/jwt")
		//middleware//
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod:"HS256",
		SigningKey: []byte("secret"),
	}))


	jwtGroup.GET("/login",jwtLogin)

	//router
	e.POST("/register", regUser)
	e.GET("/login_tkn", loginUser)
	
	e.Logger.Fatal(e.Start(":8080"))
}
