package e2e_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxiOtero6/go-auth-rest/model"
	"github.com/MaxiOtero6/go-auth-rest/repository"
	"github.com/MaxiOtero6/go-auth-rest/service"
	"github.com/MaxiOtero6/go-auth-rest/test"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	jwtService := service.NewJWTService(nil)

	testUser := model.User{
		Username: "test", Email: "test@test.com",
	}

	token, err := jwtService.Sign(testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	t.Run("It should retrieve an unauthorized status code if not auth header is provided", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Unauthorized user", data["title"])
		assert.Equal(t, "The user isn't authorized because no Authorization header is provided", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/users/", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve an empty array if the table is empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}
		assert.Empty(t, data, "If the users table is empty, it should return an empty array")
	})

	t.Run("It should retrieve all the users in the table", func(t *testing.T) {
		defer test.ClearUsers()
		w := httptest.NewRecorder()
		repository := repository.UserRepository{}

		repository.CreateUser(&model.BaseUser{Username: "test1", Email: "test1@test.com", Password: "test1"})
		repository.CreateUser(&model.BaseUser{Username: "test2", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "test3", Email: "test3@test.com", Password: "test3"})

		req, _ := http.NewRequest(http.MethodGet, "/users/", nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data []model.User
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a []model.User parseable string, ", err)
		}

		assert.NotEmpty(t, data, "If the users table is not empty, it should return a non-empty array")
		assert.Len(t, data, 3, "The length of the array should be 3")
		assert.Equal(t, model.User{Username: "test3", Email: "test3@test.com"}, data[0])
		assert.Equal(t, model.User{Username: "test2", Email: "test2@test.com"}, data[1])
		assert.Equal(t, model.User{Username: "test1", Email: "test1@test.com"}, data[2])
	})
}

func TestGetUser(t *testing.T) {
	jwtService := service.NewJWTService(nil)

	testUser := model.User{
		Username: "test", Email: "test@test.com",
	}

	token, err := jwtService.Sign(testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	t.Run("It should retrieve an unauthorized status code if not auth header is provided", func(t *testing.T) {
		username := "test"
		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Unauthorized user", data["title"])
		assert.Equal(t, "The user isn't authorized because no Authorization header is provided", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/users/"+username, data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a not found status if the user with username 'test' is not found", func(t *testing.T) {
		username := "test"

		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")
		log.Println(w.Body.String())
		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "User not found", data["title"])
		assert.Equal(t, "The user with the username "+username+" was not found", data["detail"])
		assert.Equal(t, float64(404), data["status"])
		assert.Equal(t, "/users/"+username, data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve only the user with username 'test'", func(t *testing.T) {
		username := "test"
		w := httptest.NewRecorder()

		repository := repository.UserRepository{}

		repository.CreateUser(&model.BaseUser{Username: testUser.Username, Email: testUser.Email, Password: "test"})
		repository.CreateUser(&model.BaseUser{Username: "test2", Email: "test2@test.com", Password: "test2"})
		repository.CreateUser(&model.BaseUser{Username: "test3", Email: "test3@test.com", Password: "test3"})
		defer test.ClearUsers()

		req, _ := http.NewRequest(http.MethodGet, "/users/"+username, nil)
		req.Header.Add("Authorization", bearerToken)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var data model.User
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a model.User parseable string, ", err)
		}

		assert.Equal(t, testUser, data)
	})
}
