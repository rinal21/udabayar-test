package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/database"
	"udabayar-go-api-di/handler"
	"udabayar-go-api-di/helper"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ShopTestSuite struct {
	suite.Suite
	ShopHandler *handler.ShopHandler
	AuthHandler *handler.AuthHandler
}

type ShopInput struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
	Phone string `faker:"phone"`
}

func (s *ShopTestSuite) setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("", s.AuthHandler.Authorizer, s.ShopHandler.AllShop)
	router.GET("/:id", s.ShopHandler.GetShop)
	// shop.POST("/:id/courier", s.AuthHandler.Authorizer, s.ShopCourier.CreateShopCourier)
	router.POST("", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.CreateShop)
	router.PUT("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.UpdateShop)
	router.DELETE("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.DeleteShop)
	router.GET("/:id/address", s.ShopHandler.GetMainAddress)

	return router
}

func TestShop(t *testing.T) {
	suite.Run(t, new(ShopTestSuite))
}

func (s *ShopTestSuite) SetupTest() {
	os.Setenv("ENVIRONMENT", "testing")

	dbconfig := config.NewDatabaseConfig()
	db, _ := database.ConnectDatabase(dbconfig)

	dbmigration := database.DatabaseMigration(db)
	dbmigration.Migrate()

	mainConfig := config.NewConfig()
	jwt := helper.NewJWTHelper(mainConfig)

	bcrypt := helper.NewBcryptHelper()

	emailconfig := config.NewEmailConfig()
	email := helper.NewEmailHelper(emailconfig)

	s.ShopHandler = handler.NewShopHandler(db)
	s.AuthHandler = handler.NewAuthHandler(mainConfig, db, jwt, bcrypt, email)
}

func (s *ShopTestSuite) Test7AllShop() {
	router := s.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *ShopTestSuite) Test8GetShop() {
	router := s.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/1", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *ShopTestSuite) Test9GetMainAddress() {
	router := s.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/1/address", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

// func (s *ShopTestSuite) Test6UpdateShop() {
// 	router := s.setupRouter()

// 	UpdateShopInput := ShopInput{
// 		Name:  "Andi",
// 		Email: "tester123@gmail.com",
// 		Phone: "0987654",
// 	}

// 	params := url.Values{}
// 	params.Add("name", UpdateShopInput.Name)
// 	params.Add("email", UpdateShopInput.Email)
// 	params.Add("phone", UpdateShopInput.Phone)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("PUT", "/", strings.NewReader(params.Encode()))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
// 	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

// 	router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// }
