package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NutriPocket/UserService/database"
	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/NutriPocket/UserService/test"
	"github.com/stretchr/testify/assert"
)

func TestPostRegister(t *testing.T) {
	gormDB, err := database.GetPoolConnection()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	type RequestData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	t.Run("It should create a new user", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: "test@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		var resData map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resData)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		var data map[string]interface{} = resData["data"].(map[string]interface{})

		assert.Equal(t, bodyData.Username, data["username"], "Username should be "+bodyData.Username)
		assert.Equal(t, bodyData.Email, data["email"], "Email should be"+bodyData.Email)
		assert.NotEmpty(t, resData["token"], "Token should not be empty")

		var saved model.User

		gormDB.Raw("SELECT username, email FROM users WHERE username = ?", bodyData.Username).Scan(&saved)

		assert.Equal(t, bodyData.Username, saved.Username)
	})

	t.Run("It should retrieve a bad request status if the provided data is blank", func(t *testing.T) {
		defer test.ClearUsers()

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		assert.Equal(t, "Wrong body format", data["title"])
		assert.Equal(t, "Expected a json body with the user account data in it", data["detail"])
		assert.Equal(t, float64(http.StatusBadRequest), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])

		var saved model.User

		gormDB.Raw("SELECT username, email FROM users WHERE username = ?", "").Scan(&saved)

		assert.Equal(t, "", saved.Username)
	})

	t.Run("It should retrieve a bad request status if the provided username is blank", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "", Email: "test@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Empty username field", data["title"])
		assert.Equal(t, "The username field is required", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided username is too long", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: strings.Repeat("username", 20), Email: "test@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Invalid username field", data["title"])
		assert.Equal(t, "The username field must be less than 100 characters", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a conflict status if the provided username is already in use", func(t *testing.T) {
		defer test.ClearUsers()

		res := gormDB.Exec(
			`
				INSERT INTO users (username, email, password) 
				VALUES (?, ?, ?)
			`,
			"test", "test@test.com", "test",
		)

		if res.Error != nil {
			log.Fatalf("An error ocurred when inserting a user: %v\n", res.Error)
		}

		bodyData := RequestData{Username: "test", Email: "test2@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code, "Status code should be 409")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Username or email already in use", data["title"])
		assert.Equal(t, "The provided username or email are already in use, try something else", data["detail"])
		assert.Equal(t, float64(409), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided password is blank", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: "test@test.com", Password: ""}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Empty password field", data["title"])
		assert.Equal(t, "The password field is required", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided password is too long", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: "test@test.com", Password: strings.Repeat("password", 20)}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Invalid password field", data["title"])
		assert.Equal(t, "The password field must be less than 100 characters", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided email is blank", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: "", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Empty email field", data["title"])
		assert.Equal(t, "The email field is required", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided email is invalid", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: "testtest.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Invalid email field", data["title"])
		assert.Equal(t, "The email field must be a valid email address", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided email is too long", func(t *testing.T) {
		defer test.ClearUsers()

		bodyData := RequestData{Username: "test", Email: strings.Repeat("test@test.com", 20), Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Invalid email field", data["title"])
		assert.Equal(t, "The email field must be less than 100 characters", data["detail"])
		assert.Equal(t, float64(400), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a conflict status if the provided email is already in use", func(t *testing.T) {
		defer test.ClearUsers()

		res := gormDB.Exec(
			`
				INSERT INTO users (username, email, password) 
				VALUES (?, ?, ?)
			`,
			"test", "test@test.com", "test",
		)

		if res.Error != nil {
			log.Fatalf("An error ocurred when inserting a user: %v\n", res.Error)
		}

		bodyData := RequestData{Username: "test2", Email: "test@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/register", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code, "Status code should be 409")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Username or email already in use", data["title"])
		assert.Equal(t, "The provided username or email are already in use, try something else", data["detail"])
		assert.Equal(t, float64(409), data["status"])
		assert.Equal(t, "/auth/register", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})
}

func TestPostLogin(t *testing.T) {
	gormDB, err := database.GetPoolConnection()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	defer test.ClearUsers()

	type RequestData struct {
		EmailOrUsername string `json:"emailOrUsername"`
		Password        string `json:"password"`
	}

	service, err := service.NewUserService(nil)
	if err != nil {
		log.Fatalf("An error ocurred when creating the user service: %v\n", err)
	}

	encodedPassword := service.EncodePassword("test")

	res := gormDB.Exec(
		`
			INSERT INTO users (username, email, password) 
			VALUES (?, ?, ?)
		`,
		"test", "test@test.com", encodedPassword,
	)

	if res.Error != nil {
		log.Fatalf("An error ocurred when inserting a user: %v\n", res.Error)
	}

	t.Run("It should login using username and password", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: "test", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var resData map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resData)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		var data map[string]interface{} = resData["data"].(map[string]interface{})

		assert.Equal(t, "test", data["username"], "Username should be test")
		assert.Equal(t, "test@test.com", data["email"], "Email should be test@test.com")
		assert.NotEmpty(t, resData["token"], "Token should not be empty")
	})

	t.Run("It should login using email and password", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: "test@test.com", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var resData map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resData)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		var data map[string]interface{} = resData["data"].(map[string]interface{})

		assert.Equal(t, "test", data["username"], "Username should be test")
		assert.Equal(t, "test@test.com", data["email"], "Email should be test@test.com")
		assert.NotEmpty(t, resData["token"], "Token should not be empty")
	})

	t.Run("It should retrieve a bad request status if the provided data is blank", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/auth/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		assert.Equal(t, "Wrong body format", data["title"])
		assert.Equal(t, "Expected a json body with the user credentials in it", data["detail"])
		assert.Equal(t, float64(http.StatusBadRequest), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])

		var saved model.User

		gormDB.Raw("SELECT username, email FROM users WHERE username = ?", "").Scan(&saved)

		assert.Equal(t, "", saved.Username)
	})

	t.Run("It should retrieve a bad request status if the provided username or email is blank", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: "", Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Credentials don't match", data["title"])
		assert.Equal(t, "User identification or password are wrong, please try again", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided username is too long", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: strings.Repeat("username", 20), Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Credentials don't match", data["title"])
		assert.Equal(t, "User identification or password are wrong, please try again", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided password is blank", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: "test@test.com", Password: ""}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Credentials don't match", data["title"])
		assert.Equal(t, "User identification or password are wrong, please try again", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided password is too long", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: "test@test.com", Password: strings.Repeat("password", 20)}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Credentials don't match", data["title"])
		assert.Equal(t, "User identification or password are wrong, please try again", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})

	t.Run("It should retrieve a bad request status if the provided email is too long", func(t *testing.T) {
		bodyData := RequestData{EmailOrUsername: strings.Repeat("test@test.com", 20), Password: "test"}
		jsonData, err := json.Marshal(bodyData)

		if err != nil {
			log.Fatal("The data is not a JSON parseable object, ", err)
		}

		body := bytes.NewBuffer(jsonData)

		req, _ := http.NewRequest(http.MethodPost, "/auth/login", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Credentials don't match", data["title"])
		assert.Equal(t, "User identification or password are wrong, please try again", data["detail"])
		assert.Equal(t, float64(401), data["status"])
		assert.Equal(t, "/auth/login", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})
}

func TestPostLogout(t *testing.T) {
	gormDB, err := database.GetPoolConnection()
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	type RequestData struct {
		Token string `json:"token"`
	}

	type BlacklistedToken struct {
		Signature string `json:"signature"`
	}

	jwtService, err := service.NewJWTService(nil)
	if err != nil {
		log.Fatalf("An error ocurred when creating the JWT service: %v\n", err)
	}

	testUser := model.User{
		Username: "test", Email: "test@test.com",
	}

	token, err := jwtService.Sign(testUser)

	if err != nil {
		log.Fatalf("An error ocurred when signing testUser: %v\n", err)
	}

	lastDotIndex := strings.LastIndex(token, ".")
	signature := token[lastDotIndex:]

	bodyData := RequestData{Token: token}
	jsonData, err := json.Marshal(bodyData)

	if err != nil {
		log.Fatal("The data is not a JSON parseable object, ", err)
	}

	t.Run("It should blacklist the token signature", func(t *testing.T) {
		defer test.ClearBlacklist()
		body := bytes.NewBuffer(jsonData)
		req, _ := http.NewRequest(http.MethodPost, "/auth/logout", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code, "Status code should be 204")

		var saved BlacklistedToken

		gormDB.Raw("SELECT signature FROM jwt_blacklist WHERE signature = ?", signature).Scan(&saved)

		assert.Equal(t, signature, saved.Signature)
	})

	t.Run("It should retrieve a bad request status if the provided token is blank", func(t *testing.T) {
		defer test.ClearBlacklist()

		req, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var data map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		assert.Equal(t, "Wrong body format", data["title"])
		assert.Equal(t, "Expected a json body with the key 'token' in it", data["detail"])
		assert.Equal(t, float64(http.StatusBadRequest), data["status"])
		assert.Equal(t, "/auth/logout", data["instance"])
		assert.Equal(t, "about:blank", data["type"])

		var saved BlacklistedToken

		gormDB.Raw("SELECT signature FROM jwt_blacklist WHERE signature = ?", signature).Scan(&saved)

		assert.Equal(t, "", saved.Signature)
	})

	t.Run("It should retrieve a conflict status if the provided token is already blacklisted", func(t *testing.T) {
		defer test.ClearBlacklist()
		body := bytes.NewBuffer(jsonData)
		gormDB.Exec(
			`
				INSERT INTO jwt_blacklist (signature, expires_at)
				VALUES (?, NOW() + INTERVAL 1 DAY)
			`,
			signature,
		)

		req, _ := http.NewRequest(http.MethodPost, "/auth/logout", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code, "Status code should be 409")

		var data map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &data)
		if err != nil {
			log.Fatal("The response body is not a JSON parseable string, ", err)
		}

		assert.Equal(t, "Token no longer used", data["title"])
		assert.Equal(t, "The provided token is no longer in use", data["detail"])
		assert.Equal(t, float64(409), data["status"])
		assert.Equal(t, "/auth/logout", data["instance"])
		assert.Equal(t, "about:blank", data["type"])
	})
}
