package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret")

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

//Claims struct
type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	jwt.StandardClaims
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
func regUser(c echo.Context) error {
	reg := new(Register)
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&reg)
	if err != nil {
		log.Printf("failed processing request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("Register : %#v", reg)
	hash, err := bcrypt.GenerateFromPassword([]byte(reg.PasswordHash), bcrypt.DefaultCost)
	reg.PasswordHash = string(hash)

	if err != nil {
		log.Println(err)
	}
	err = db.Table("registers").Create(&reg).Error
	if err != nil {
		log.Println(err)
	}
	return c.JSON(http.StatusCreated, &reg)
}
func loginUser(c echo.Context) error {
	lgp := new(Register)
	email := c.QueryParam("email")
	password := c.QueryParam("password")
	err := db.Where("email = ? ", email).First(&lgp).Error
	if err != nil {
		log.Println("db error")
	}
	// fmt.Printf("hash value: %s \n plain password: %s \n", lgp.PasswordHash, password)
	err1 := bcrypt.CompareHashAndPassword([]byte(lgp.PasswordHash), []byte(password))

	if err1 != nil {
		log.Println(err1)
	} else {

		//create tooken

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Name:  lgp.Name,
			Email: lgp.Email,
			Phone: lgp.Phone,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tkn, err := tk.SignedString(jwtKey)
		if err != nil {
			log.Printf("failed processing request: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		cookie := new(http.Cookie)
		cookie.Name = "tooken"
		cookie.Value = tkn
		c.SetCookie(cookie)
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
func whoAmI(c echo.Context) error {
	cookie, err := c.Cookie("tooken")
	if err != nil {
		return echo.ErrUnauthorized
	}
	tknStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(tooken *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {

		if err == jwt.ErrSignatureInvalid {
			return echo.ErrUnauthorized
		}

		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !tkn.Valid {
		return echo.ErrUnauthorized
	}
	user := new(Register)
	err1 := db.Where("email = ? ", claims.Email).First(&user).Error
	if err1 != nil {
		log.Println("db error")
	}
	id := strconv.FormatInt(user.ID, 10)
	return c.JSON(http.StatusOK, map[string]string{

		"id":    id,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
	})
}
func main() {
	initDb()
	e := echo.New()
	wg := new(sync.WaitGroup)

	// add two goroutines to `wg` WaitGroup
	wg.Add(3)
	//jwt group//
	jwtGroup := e.Group("api/v1/user/")
	//middleware//
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte("secret"),
	}))

	jwtGroup.POST("/whoAmI", whoAmI)

	//router
	e.POST("api/v1/user/register", regUser)
	e.POST("api/v1/user/login_tkn", loginUser)

	e.GET("api/v1/auth/users", User)
	e.GET("api/v1/auth/userid", userByID)
	go func() {

		e.Logger.Fatal(e.Start(":8080"))
		wg.Done()
	}()
	go func() {

		e.Logger.Fatal(e.Start(":8081"))
		wg.Done()
	}()
	go func() {

		e.Logger.Fatal(e.Start(":8082"))
		wg.Done()
	}()

	wg.Wait()

}
