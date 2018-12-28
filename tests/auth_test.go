package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/database"
	"udabayar-go-api-di/handler"
	"udabayar-go-api-di/helper"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	AuthHandler *handler.AuthHandler
}

type AuthInput struct {
	Name     string `faker:"name"`
	Email    string `faker:"email"`
	Password string `faker:"password"`
}

type ChangePasswordInput struct {
	CurrentPassword string `faker:"current_password"`
	NewPassword     string `faker:"new_password"`
}

func (s *AuthTestSuite) setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/register", s.AuthHandler.Register)
	router.POST("/login", s.AuthHandler.Login)
	// router.GET("/activation", s.AuthHandler.ActivateEmail)
	router.POST("/validatetoken", s.AuthHandler.ValidateToken)
	router.POST("/change_password", s.AuthHandler.Authorizer, s.AuthHandler.ChangePassword)

	return router
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupTest() {
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

	s.AuthHandler = handler.NewAuthHandler(mainConfig, db, jwt, bcrypt, email)
}

// func (s *AuthTestSuite) Test1AuthRegister() {
// 	router := s.setupRouter()

// 	registerInput := AuthInput{
// 		Name:     "si Tester",
// 		Email:    "tester123@gmail.com",
// 		Password: "tester99123",
// 	}

// 	params := url.Values{}
// 	params.Add("name", registerInput.Name)
// 	params.Add("email", registerInput.Email)
// 	params.Add("password", registerInput.Password)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/register", strings.NewReader(params.Encode()))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

// 	router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// }

func (s *AuthTestSuite) Test2AuthLogin() {
	router := s.setupRouter()

	loginInput := AuthInput{
		Email:    "tester123@gmail.com",
		Password: "tester99123",
	}

	params := url.Values{}
	params.Add("email", loginInput.Email)
	params.Add("password", loginInput.Password)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

// func (s *AuthTestSuite) Test3ValidateToken() {
// 	router := s.setupRouter()

// 	mainConfig := config.NewConfig()
// 	jwt := helper.NewJWTHelper(mainConfig)
// 	generateToken, _ := jwt.SignJWT(1, "access")

// 	params := url.Values{}
// 	params.Add("token", generateToken)

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("POST", "/validatetoken", strings.NewReader(params.Encode()))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

// 	router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// }

func (s *AuthTestSuite) Test4ChangePassword() {
	router := s.setupRouter()

	Value := ChangePasswordInput{
		CurrentPassword: "tester99123",
		NewPassword:     "tester99123",
	}

	params := url.Values{}
	params.Add("current_password", Value.CurrentPassword)
	params.Add("new_password", Value.NewPassword)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/change_password", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

// func (s *AuthTestSuite) TestActivationEmail() {
// 	router := s.setupRouter()

// 	mainConfig := config.NewConfig()
// 	jwt := helper.NewJWTHelper(mainConfig)
// 	generateToken, _ := jwt.SignJWT(1, "emailconfirmation")

// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/activation?token="+generateToken, nil)

// 	router.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// }
