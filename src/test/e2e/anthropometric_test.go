package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/test"
	"github.com/stretchr/testify/assert"
)

func unmarshallAnthropometricData(t *testing.T, body []byte) model.AnthropometricData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.AnthropometricData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func TestPutUserAnthropometrics(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/anthropometrics/", userId)

	t.Run("PUT /users/:userId/anthropometrics - Create Anthropometric Data with all fields", func(t *testing.T) {
		defer test.ClearAllData()

		muscleMass := float32(30.2)
		fatMass := float32(20.1)
		boneMass := float32(10.0)

		payload := model.AnthropometricData{
			UserID:     userId,
			Weight:     70.5,
			MuscleMass: &muscleMass,
			FatMass:    &fatMass,
			BoneMass:   &boneMass,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Create Anthropometric Data without optional fields", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.AnthropometricData{
			UserID: userId,
			Weight: 70.5,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Create Anthropometric Data with muscle mass field", func(t *testing.T) {
		defer test.ClearAllData()

		muscleMass := float32(30.2)

		payload := model.AnthropometricData{
			UserID:     userId,
			Weight:     70.5,
			MuscleMass: &muscleMass,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Create Anthropometric Data with fat mass field", func(t *testing.T) {
		defer test.ClearAllData()

		fatMass := float32(20.1)

		payload := model.AnthropometricData{
			UserID:  userId,
			Weight:  70.5,
			FatMass: &fatMass,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Create Anthropometric Data with bone mass field", func(t *testing.T) {
		defer test.ClearAllData()

		boneMass := float32(10.0)

		payload := model.AnthropometricData{
			UserID:   userId,
			Weight:   70.5,
			BoneMass: &boneMass,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Update Anthropometric Data", func(t *testing.T) {
		defer test.ClearAllData()

		{
			muscleMass := float32(30.2)
			fatMass := float32(20.1)
			boneMass := float32(10.0)
			payload := model.AnthropometricData{
				UserID:     userId,
				Weight:     70.5,
				MuscleMass: &muscleMass,
				FatMass:    &fatMass,
				BoneMass:   &boneMass,
			}

			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		muscleMass := float32(32.0)
		fatMass := float32(18.0)
		boneMass := float32(11.0)
		newPayload := model.AnthropometricData{
			UserID:     userId,
			Weight:     75.0,
			MuscleMass: &muscleMass,
			FatMass:    &fatMass,
			BoneMass:   &boneMass,
		}

		body, _ := json.Marshal(newPayload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, newPayload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Update Anthropometric Data with optional fields that were null in the original entry", func(t *testing.T) {
		defer test.ClearAllData()

		{
			payload := model.AnthropometricData{
				UserID: userId,
				Weight: 70.5,
			}

			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		muscleMass := float32(32.0)
		fatMass := float32(18.0)
		boneMass := float32(11.0)
		newPayload := model.AnthropometricData{
			UserID:     userId,
			Weight:     75.0,
			MuscleMass: &muscleMass,
			FatMass:    &fatMass,
			BoneMass:   &boneMass,
		}

		body, _ := json.Marshal(newPayload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, newPayload, actual)
	})

	t.Run("PUT /users/:userId/anthropometrics - Missing weight should raise validation error", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.AnthropometricData{
			UserID: userId,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		expected := model.ErrorRfc9457{
			Title:    "Invalid anthropometric user data",
			Detail:   `The anthropometric user data is invalid, Key: 'AnthropometricData.Weight' Error:Field validation for 'Weight' failed on the 'required' tag`,
			Status:   http.StatusBadRequest,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})

	t.Run("PUT /users/:userId/anthropometrics - No token should raise Authentication Error", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.AnthropometricData{
			UserID: userId,
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
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

func TestGetUserAnthropometrics(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/anthropometrics/", userId)

	t.Run("GET /users/:userId/anthropometrics - Retrieve All Anthropometric Data should be empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		var response map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Empty(t, response["data"], "Data should be empty")
	})

	t.Run("GET /users/:userId/anthropometrics?date=<date> - Retrieve Data by Date found data", func(t *testing.T) {
		defer test.ClearAllData()

		muscleMass := float32(30.2)
		fatMass := float32(20.1)
		boneMass := float32(10.0)
		payload := model.AnthropometricData{
			UserID:     userId,
			Weight:     70.5,
			MuscleMass: &muscleMass,
			FatMass:    &fatMass,
			BoneMass:   &boneMass,
		}

		{
			putURL := fmt.Sprintf("/users/%s/anthropometrics/", userId)
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, putURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		date := time.Now().Format("2006-01-02")
		url := fmt.Sprintf("%s?date=%s", baseURL, date)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallAnthropometricData(t, w.Body.Bytes())

		actual.CreatedAt = ""

		assert.Equal(t, payload, actual)
	})

	t.Run("GET /users/:userId/anthropometrics?date=<date> - Retrieve Data by Date data not found", func(t *testing.T) {
		date := "2023-01-01"
		url := fmt.Sprintf("%s?date=%s", baseURL, date)

		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		expected := model.ErrorRfc9457{
			Title:    "Anthropometric data not found",
			Detail:   "No anthropometric data not found for user " + userId + " on date " + date,
			Status:   http.StatusNotFound,
			Type:     "about:blank",
			Instance: url,
		}

		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})

	t.Run("GET /users/:userId/anthropometrics - No token should raise Authentication Error", func(t *testing.T) {
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
