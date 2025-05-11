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

func unmarshallObjectiveData(t *testing.T, body []byte) model.ObjectiveData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.ObjectiveData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func TestPutUserObjective(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/objectives/", userId)

	t.Run("PUT /users/:userId/objectives - Create Objective", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID:     userId,
				Weight:     70.0,
				MuscleMass: floatPtr(30.0),
				FatMass:    floatPtr(20.0),
			},
			Deadline: "2025-12-31",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallObjectiveData(t, w.Body.Bytes())

		actual.Deadline = strings.Split(actual.Deadline, "T")[0]
		actual.CreatedAt = ""

		assert.Equal(t, payload, actual, "Updated objective data should match")
	})

	t.Run("PUT /users/:userId/objectives - Past deadlines should raise validation error", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID:     userId,
				Weight:     70.0,
				MuscleMass: floatPtr(30.0),
				FatMass:    floatPtr(20.0),
			},
			Deadline: "2024-12-31",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		expected := model.ErrorRfc9457{
			Title:    "Invalid deadline",
			Detail:   `The deadline must be in the future`,
			Status:   http.StatusBadRequest,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})

	t.Run("PUT /users/:userId/objectives - No deadline should raise validation error", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID:     userId,
				Weight:     70.0,
				MuscleMass: floatPtr(30.0),
				FatMass:    floatPtr(20.0),
			},
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		expected := model.ErrorRfc9457{
			Title:    "Invalid user objective data",
			Detail:   `The user objective data is invalid, Key: 'ObjectiveData.Deadline' Error:Field validation for 'Deadline' failed on the 'required' tag`,
			Status:   http.StatusBadRequest,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})

	t.Run("PUT /users/:userId/objectives - Update Objective", func(t *testing.T) {
		defer test.ClearAllData()

		// Initial data
		initialPayload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID: userId,
				Weight: 70.0,
			},
			Deadline: "2025-12-31",
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
		updatedPayload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID:     userId,
				Weight:     75.0,
				MuscleMass: floatPtr(32.0),
				FatMass:    floatPtr(18.0),
			},
			Deadline: "2026-01-01",
		}

		body, _ := json.Marshal(updatedPayload)
		req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallObjectiveData(t, w.Body.Bytes())

		actual.Deadline = strings.Split(actual.Deadline, "T")[0]
		actual.CreatedAt = ""

		assert.Equal(t, updatedPayload, actual, "Updated objective data should match")
	})
}

func TestGetUserObjective(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/objectives/", userId)

	t.Run("GET /users/:userId/objectives - Retrieve Objective", func(t *testing.T) {
		defer test.ClearAllData()

		// Create objective
		payload := model.ObjectiveData{
			AnthropometricData: model.AnthropometricData{
				UserID:     userId,
				Weight:     70.0,
				MuscleMass: floatPtr(30.0),
				FatMass:    floatPtr(20.0),
			},
			Deadline: "2025-12-31",
		}

		{
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Retrieve objective
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallObjectiveData(t, w.Body.Bytes())

		actual.Deadline = strings.Split(actual.Deadline, "T")[0]
		actual.CreatedAt = ""

		assert.Equal(t, payload, actual, "Updated objective data should match")
	})

	t.Run("GET /users/:userId/objectives - No objetive should raise not found error", func(t *testing.T) {
		// Retrieve objective
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

		expected := model.ErrorRfc9457{
			Title:    "Objective data not found",
			Detail:   "No objective data not found for user " + userId,
			Status:   http.StatusNotFound,
			Type:     "about:blank",
			Instance: baseURL,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})

	t.Run("GET /users/:userId/objectives - No token should raise Authentication Error", func(t *testing.T) {
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

func floatPtr(f float32) *float32 {
	return &f
}
