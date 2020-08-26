package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

//Register struct//
type Register struct {
	// gorm.Model
	ID       int    `json:"id" gorm:"primary_key;AUTO_INCREMENT;"`
	Email    string `json:"email,omitempty" gorm:"type:varchar(100);" `
	Name     string `json:"name,omitempty" gorm:"type:varchar(50);" `
	Phone    string `json:"phone,omitempty" gorm:"type:varchar(20);unique;" `
	Password string `json:"password,omitempty" gorm:"type:varchar(300);" `
}

var db *gorm.DB

func initDb() {
	var err error
	db, err = gorm.Open("mysql", "root:itsshawn@007@@tcp(localhost:3306)/signs?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect Database")
	}
	// db.Exec("CREATE DATABASE signs")

	db.Exec("use signs")
	db.AutoMigrate(&Register{})
}
func regUser(c echo.Context) error {
	reg := new(Register)
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reg)
	if err != nil {
		log.Printf("failed processing request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Register : %#v", reg)
	hash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	reg.Password = string(hash)
	if err != nil {
		log.Println(err)
	}
	err = db.Table("registers").Create(&reg).Error
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, &reg)
	// return c.String(http.StatusOK,"You are registered")
}
func loginUser(c echo.Context) error {
	lgp := new(Register)
	Email := c.QueryParam("email")
	Password := c.QueryParam("password")

	err := db.Where("email = ? AND password= ? ", Email, Password).First(&lgp).Error
	err1 := bcrypt.CompareHashAndPassword([]byte(lgp.Password), []byte(Password))
	if err1 != nil {
		log.Println(err1)
	}
	lgp.Password = Password
	if err != nil {
		log.Println(err)
	} else {

		//create tooken
		tk := jwt.New(jwt.SigningMethodHS256)
		claims := tk.Claims.(jwt.MapClaims)
		claims["name"] = lgp.Name
		claims["email"] = lgp.Email
		claims["phone"] = lgp.Phone
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
func jwtLogin(c echo.Context) error {

	return c.String(http.StatusOK, "You are Logged in sucessfully !")
}
func main() {
	initDb()
	e := echo.New()

	//jwt group//
	jwtGroup := e.Group("api/v1/user/jwt")
	//middleware//
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte("secret"),
	}))

	jwtGroup.POST("/login", jwtLogin)

	//router
	e.POST("api/v1/user/register", regUser)
	e.POST("api/v1/user/login_tkn", loginUser)

	e.GET("api/v1/auth/user", User)
	e.GET("api/v1/auth/userid", userByID)

	e.Logger.Fatal(e.Start(":8080"))
}
