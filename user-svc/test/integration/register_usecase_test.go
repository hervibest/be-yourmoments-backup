package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hervibest/be-yourmoments-backup/user-svc/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := webServer(ctx); err != nil {
			log.Fatalf("Failed to start web server: %v", err)
		}
	}()

	if err := waitForServerReady(); err != nil {
		log.Fatalf("Server not ready in time: %v", err)
	}

	code := m.Run()

	cancel()

	os.Exit(code)
}

func waitForServerReady() error {
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://127.0.0.1:8003/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("server not ready after %d retries", maxRetries)
}

func TestRegisterByPhoneNumber(t *testing.T) {
	// Clear existing users and user profiles
	ClearUserProfiles()
	ClearUsers()
	requestBody := model.RegisterByPhoneRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		PhoneNumber:  "085228561067",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/phone", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.Equal(t, requestBody.PhoneNumber, *responseBody.Data.PhoneNumber)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestDupicatePhoneRegisterByPhoneNumber(t *testing.T) {
	TestRegisterByPhoneNumber(t)
	requestBody := model.RegisterByPhoneRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		PhoneNumber:  "085228561067",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/phone", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "Phone number has already been taken", responseBody.Message)
}

func TestDupicateUsernameRegisterByPhoneNumber(t *testing.T) {
	TestRegisterByPhoneNumber(t)
	requestBody := model.RegisterByPhoneRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		PhoneNumber:  "085228561063",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/phone", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "Username has already been taken", responseBody.Message)
}

func TestInvalidRegisterByPhoneNumber(t *testing.T) {
	// Clear existing users and user profiles
	ClearUserProfiles()
	ClearUsers()
	requestBody := model.RegisterByPhoneRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		PhoneNumber:  "085228561063333333337",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/phone", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ValidationErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnprocessableEntity, response.StatusCode)
	assert.Equal(t, "Validation error", responseBody.Message)
	assert.Equal(t, 1, len(responseBody.Errors))
}

func TestRegisterByEmail(t *testing.T) {
	ClearEmailVerifications()
	ClearUserProfiles()
	ClearUsers()
	requestBody := model.RegisterByEmailRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		Email:        "hervipro@gmail.com",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/email", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, 100000)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.Equal(t, requestBody.Email, *responseBody.Data.Email)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestDuplicateEmailRegisterByEmail(t *testing.T) {
	TestRegisterByEmail(t)
	requestBody := model.RegisterByEmailRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		Email:        "hervipro@gmail.com",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/email", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, 100000)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "Email has already been taken", responseBody.Message)
}

func TestDuplicateUsernameRegisterByEmail(t *testing.T) {
	TestRegisterByEmail(t)
	requestBody := model.RegisterByEmailRequest{
		Username:     "hervi",
		Password:     "hervi12345!",
		Email:        "hervipro@gmail.com",
		BirthDateStr: "2001-10-05",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/user/register/email", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, 100000)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse)
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "Username has already been taken", responseBody.Message)
}
