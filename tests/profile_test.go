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

type ProfileTestSuite struct {
	suite.Suite
	ProfileHandler *handler.ProfileHandler
}

type ProfileInput struct {
	Name  string `faker:"name"`
	Email string `faker:"email"`
	Phone string `faker:"phone"`
}

func (s *ProfileTestSuite) setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("", s.ProfileHandler.GetProfile)
	router.PUT("", s.ProfileHandler.UpdateProfile)
	router.PUT("/image", s.ProfileHandler.UpdateProfileImage)

	return router
}

func TestProfile(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (s *ProfileTestSuite) SetupTest() {
	os.Setenv("ENVIRONMENT", "testing")

	dbconfig := config.NewDatabaseConfig()
	db, _ := database.ConnectDatabase(dbconfig)

	dbmigration := database.DatabaseMigration(db)
	dbmigration.Migrate()

	mainConfig := config.NewConfig()

	upload := helper.NewUploadHelper()

	s.ProfileHandler = handler.NewProfileHandler(mainConfig, db, upload)
}

func (s *ProfileTestSuite) Test5Profile() {
	router := s.setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}

func (s *ProfileTestSuite) Test6UpdateProfile() {
	router := s.setupRouter()

	UpdateProfileInput := ProfileInput{
		Name:  "Andi",
		Email: "tester123@gmail.com",
		Phone: "0987654",
	}

	params := url.Values{}
	params.Add("name", UpdateProfileInput.Name)
	params.Add("email", UpdateProfileInput.Email)
	params.Add("phone", UpdateProfileInput.Phone)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	req.Header.Add("Authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEzLCJTY29wZSI6ImFwaTphY2Nlc3MiLCJleHAiOjE1NDU5ODMzMjV9.Ld9oxxYN7mvU_WY9xZKfWCeHh1h57OIOMADo8fnndSA")

	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}
