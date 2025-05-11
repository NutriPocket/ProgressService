package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/test"
	"github.com/stretchr/testify/assert"
)

func unmarshallFixedUserData(t *testing.T, body []byte) model.FixedUserData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.FixedUserData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func unmarshallBaseFixedUserData(t *testing.T, body []byte) model.BaseFixedUserData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.BaseFixedUserData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func TestPutUserFixedData(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/fixedData/", userId)

	t.Run("PUT /users/:userId/fixedData - Create Fixed Data", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.BaseFixedUserData{
			UserID:   userId,
			Height:   180,
			Birthday: "1990-01-01",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallFixedUserData(t, w.Body.Bytes())

		assert.Equal(t, payload.UserID, actual.UserID)
		assert.Equal(t, payload.Height, actual.Height)
		assert.NotEqual(t, 0, actual.Age, "Age should be calculated")
	})

	t.Run("PUT /users/:userId/fixedData - Update Fixed Data", func(t *testing.T) {
		defer test.ClearAllData()

		// Initial data
		initialPayload := model.BaseFixedUserData{
			UserID:   userId,
			Height:   180,
			Birthday: "1990-01-01",
		}

		{
			body, _ := json.Marshal(initialPayload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Updated data
		updatedPayload := model.BaseFixedUserData{
			UserID:   userId,
			Height:   185,
			Birthday: "1991-01-01",
		}

		body, _ := json.Marshal(updatedPayload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallFixedUserData(t, w.Body.Bytes())

		assert.Equal(t, updatedPayload.UserID, actual.UserID)
		assert.Equal(t, updatedPayload.Height, actual.Height)
		assert.NotEqual(t, 0, actual.Age, "Age should be calculated")
	})

	t.Run("GET /users/:userId/fixedData - No token should raise Authentication Error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, baseURL, nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		expected := model.ErrorRfc9457{
			Title:    "Unauthorized user",
			Detail:   `The user isn't authorized because no Authorization header is provided`,
			Status:   http.StatusUnauthorized,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})
}

func TestGetUserFixedData(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/fixedData/", userId)

	t.Run("GET /users/:userId/fixedData - Retrieve Fixed Data", func(t *testing.T) {
		defer test.ClearAllData()

		// Create fixed data
		payload := model.BaseFixedUserData{
			UserID:   userId,
			Height:   180,
			Birthday: "1990-01-01",
		}

		{
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Retrieve fixed data
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallFixedUserData(t, w.Body.Bytes())

		assert.Equal(t, payload.UserID, actual.UserID)
		assert.Equal(t, payload.Height, actual.Height)
		assert.NotEqual(t, 0, actual.Age, "Age should be calculated")
	})

	t.Run("GET /users/:userId/fixedData?base=true - Retrieve Base Fixed Data", func(t *testing.T) {
		defer test.ClearAllData()

		// Create fixed data
		payload := model.BaseFixedUserData{
			UserID:   userId,
			Height:   180,
			Birthday: "1990-01-01",
		}

		{
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		url := fmt.Sprintf("%s?base=true", baseURL)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallBaseFixedUserData(t, w.Body.Bytes())

		actual.Birthday = strings.Split(actual.Birthday, "T")[0]

		assert.Equal(t, payload, actual)
	})

	t.Run("GET /users/:userId/fixedData - No token should raise Authentication Error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		expected := model.ErrorRfc9457{
			Title:    "Unauthorized user",
			Detail:   `The user isn't authorized because no Authorization header is provided`,
			Status:   http.StatusUnauthorized,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})
}
