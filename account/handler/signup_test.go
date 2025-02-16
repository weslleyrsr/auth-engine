package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/weslleyrsr/auth-engine/account/model"
	"github.com/weslleyrsr/auth-engine/account/model/apperrors"
	mocks2 "github.com/weslleyrsr/auth-engine/account/model/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Email and Password Required", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks2.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks2.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob",
			"password": "supersecret1234",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks2.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"password": "sup",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Password too long", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks2.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"password": "asbdfoagbfouygasdofgoasdgfaugfp",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Error returned from UserService", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Password: "Password123",
		}

		// We just want this to show that it's not called in this case
		mockUserService := new(mocks2.MockUserService)
		mockUserService.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*model.User")).Return(apperrors.NewConflict("User Already Exists", u.Email))

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:      router,
			UserService: mockUserService,
		})

		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 409, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful Token Creation", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Password: "Password123",
		}

		mockTokenResp := &model.TokenPair{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
		}

		mockUserService := new(mocks2.MockUserService)
		mockTokenService := new(mocks2.MockTokenService)

		// Mock UserService response
		mockUserService.
			On("Signup",
				mock.AnythingOfType("*gin.Context"),
				u,
			).
			Return(nil)

		// Mock TokenService response
		mockTokenService.
			On("NewPairFromUser",
				mock.AnythingOfType("*gin.Context"),
				u,
				mock.AnythingOfType("string"),
			).
			Return(mockTokenResp, nil)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// create a request body with email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		expectedRespBody, err := json.Marshal(gin.H{
			"tokens": mockTokenResp,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, expectedRespBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})

	t.Run("Failed Token Creation", func(t *testing.T) {
		u := &model.User{
			Email:    "bob@bob.com",
			Password: "Password123",
		}

		mockErrorResponse := apperrors.NewInternal()

		mockUserService := new(mocks2.MockUserService)
		mockTokenService := new(mocks2.MockTokenService)

		// Mock UserService response
		mockUserService.
			On("Signup",
				mock.AnythingOfType("*gin.Context"),
				u,
			).
			Return(nil)

		// Mock TokenService response
		mockTokenService.
			On("NewPairFromUser",
				mock.AnythingOfType("*gin.Context"),
				u,
				mock.AnythingOfType("string"),
			).
			Return(nil, mockErrorResponse)

		// A response recorder for getting written http response
		rr := httptest.NewRecorder()

		// Don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			Router:       router,
			UserService:  mockUserService,
			TokenService: mockTokenService,
		})

		// create a request body with email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		expectedRespBody, err := json.Marshal(gin.H{
			"error": mockErrorResponse,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, expectedRespBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
		mockTokenService.AssertExpectations(t)
	})
}
