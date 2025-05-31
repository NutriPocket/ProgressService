package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NutriPocket/ProgressService/model"
	"github.com/NutriPocket/ProgressService/test"
	"github.com/stretchr/testify/assert"
)

// Helper functions to unmarshall responses
func unmarshallRoutineData(t *testing.T, body []byte) model.RoutineData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.RoutineData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func unmarshallRoutinesData(t *testing.T, body []byte) []model.RoutineData {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual []model.RoutineData
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

func unmarshallFreeSchedule(t *testing.T, body []byte) model.FreeSchedule {
	var response map[string]any
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err, "Response should be valid JSON")

	bytes, err := json.Marshal(response["data"])
	assert.NoError(t, err, "Response.data should be valid JSON")

	var actual model.FreeSchedule
	err = json.Unmarshal(bytes, &actual)
	assert.NoError(t, err, "Response.data should be valid JSON")

	return actual
}

// Test POST /users/:userId/routines/
func TestPostUserRoutine(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/routines/", userId)

	t.Run("POST /users/:userId/routines/ - Create Routine Successfully", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.RoutineDTO{
			UserID:      userId,
			Name:        "Morning Workout",
			Description: "Cardio and strength training",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 8,
				EndHour:   10,
			},
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

		actual := unmarshallRoutineData(t, w.Body.Bytes())

		assert.Equal(t, payload.UserID, actual.UserID)
		assert.Equal(t, payload.Name, actual.Name)
		assert.Equal(t, payload.Description, actual.Description)
		assert.Equal(t, payload.Day, actual.Day)
		assert.Equal(t, payload.StartHour, actual.StartHour)
		assert.Equal(t, payload.EndHour, actual.EndHour)
		assert.NotEmpty(t, actual.CreatedAt)
	})

	t.Run("POST /users/:userId/routines/ - Conflicting Schedule", func(t *testing.T) {
		defer test.ClearAllData()

		// Create first routine
		firstRoutine := model.RoutineDTO{
			UserID:      userId,
			Name:        "Morning Workout",
			Description: "Cardio and strength training",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 8,
				EndHour:   10,
			},
		}

		body, _ := json.Marshal(firstRoutine)
		req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Try to create conflicting routine
		conflictingRoutine := model.RoutineDTO{
			UserID:      userId,
			Name:        "Conflicting Workout",
			Description: "This should fail",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 9, // Overlaps with first routine
				EndHour:   11,
			},
		}

		body, _ = json.Marshal(conflictingRoutine)
		req, _ = http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code, "Status code should be 409")

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, "Routine conflict", response.Title)
		assert.Equal(t, "There is already a routine scheduled in the same time interval or subinterval", response.Detail)
	})

	t.Run("POST /users/:userId/routines/ - Unauthorized", func(t *testing.T) {
		defer test.ClearAllData()

		payload := model.RoutineDTO{
			UserID:      userId,
			Name:        "Morning Workout",
			Description: "Cardio and strength training",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 8,
				EndHour:   10,
			},
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// No authorization token
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

	t.Run("POST /users/:userId/routines/ - Invalid Data", func(t *testing.T) {
		defer test.ClearAllData()

		// Missing required fields
		payload := map[string]interface{}{
			"name": "Invalid Routine",
			// Missing day, startHour, endHour
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, "Invalid routine data", response.Title)
	})
}

// Test GET /users/:userId/routines/
func TestGetUserRoutines(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/routines/", userId)

	t.Run("GET /users/:userId/routines/ - Retrieve Routines", func(t *testing.T) {
		defer test.ClearAllData()

		// Create routines first
		routines := []model.RoutineDTO{
			{
				UserID:      userId,
				Name:        "Morning Workout",
				Description: "Cardio and strength training",
				Schedule: model.Schedule{
					Day:       "Monday",
					StartHour: 8,
					EndHour:   10,
				},
			},
			{
				UserID:      userId,
				Name:        "Evening Yoga",
				Description: "Relaxation and stretching",
				Schedule: model.Schedule{
					Day:       "Wednesday",
					StartHour: 18,
					EndHour:   19,
				},
			},
		}

		for _, routine := range routines {
			body, _ := json.Marshal(routine)
			req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Retrieve routines
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallRoutinesData(t, w.Body.Bytes())
		assert.Len(t, actual, 2, "Should return 2 routines")

		// Verify routines data
		foundRoutines := 0
		for _, r := range actual {
			if r.Name == "Morning Workout" && r.Day == "Monday" {
				foundRoutines++
				assert.Equal(t, 8, r.StartHour)
				assert.Equal(t, 10, r.EndHour)
			} else if r.Name == "Evening Yoga" && r.Day == "Wednesday" {
				foundRoutines++
				assert.Equal(t, 18, r.StartHour)
				assert.Equal(t, 19, r.EndHour)
			}
		}
		assert.Equal(t, 2, foundRoutines, "Should find both created routines")
	})

	t.Run("GET /users/:userId/routines/ - No Routines", func(t *testing.T) {
		defer test.ClearAllData()

		// Retrieve routines without creating any
		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallRoutinesData(t, w.Body.Bytes())
		assert.Empty(t, actual, "Should return empty array when no routines exist")
	})

	t.Run("GET /users/:userId/routines/ - Unauthorized", func(t *testing.T) {
		defer test.ClearAllData()

		req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
		// No authorization token
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

// Test DELETE /users/:userId/routines/
func TestDeleteUserRoutine(t *testing.T) {
	userId := testUser.ID
	baseURL := fmt.Sprintf("/users/%s/routines/", userId)

	t.Run("DELETE /users/:userId/routines/ - Delete Routine", func(t *testing.T) {
		defer test.ClearAllData()

		// Create a routine first
		routine := model.RoutineDTO{
			UserID:      userId,
			Name:        "Morning Workout",
			Description: "Cardio and strength training",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 8,
				EndHour:   10,
			},
		}

		body, _ := json.Marshal(routine)
		req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Delete the routine
		schedule := model.Schedule{
			Day:       "Monday",
			StartHour: 8,
			EndHour:   10,
		}

		body, _ = json.Marshal(schedule)
		req, _ = http.NewRequest(http.MethodDelete, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		// Verify the routine list is now empty
		actual := unmarshallRoutinesData(t, w.Body.Bytes())
		assert.Empty(t, actual, "Should return empty array after deletion")

		// Verify routine is gone by trying to get all routines
		req, _ = http.NewRequest(http.MethodGet, baseURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		actual = unmarshallRoutinesData(t, w.Body.Bytes())
		assert.Empty(t, actual, "Should return empty array when no routines exist")
	})

	t.Run("DELETE /users/:userId/routines/ - Unauthorized", func(t *testing.T) {
		defer test.ClearAllData()

		schedule := model.Schedule{
			Day:       "Monday",
			StartHour: 8,
			EndHour:   10,
		}

		body, _ := json.Marshal(schedule)
		req, _ := http.NewRequest(http.MethodDelete, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// No authorization token
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

	t.Run("DELETE /users/:userId/routines/ - Invalid Data", func(t *testing.T) {
		defer test.ClearAllData()

		// Missing required fields
		payload := map[string]interface{}{
			"day": "Monday",
			// Missing startHour, endHour
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodDelete, baseURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, "Invalid routine data", response.Title)
	})
}

// Test GET /users/freeSchedules/
func TestGetFreeSchedules(t *testing.T) {
	userId := testUser.ID
	routinesURL := fmt.Sprintf("/users/%s/routines/", userId)
	freeSchedulesURL := "/users/freeSchedules/"

	t.Run("GET /users/freeSchedules/ - Get Free Schedules", func(t *testing.T) {
		defer test.ClearAllData()

		// Create some routines first
		routines := []model.RoutineDTO{
			{
				UserID:      userId,
				Name:        "Morning Workout",
				Description: "Cardio and strength training",
				Schedule: model.Schedule{
					Day:       "Monday",
					StartHour: 8,
					EndHour:   10,
				},
			},
			{
				UserID:      userId,
				Name:        "Evening Yoga",
				Description: "Relaxation and stretching",
				Schedule: model.Schedule{
					Day:       "Wednesday",
					StartHour: 18,
					EndHour:   19,
				},
			},
		}

		for _, routine := range routines {
			body, _ := json.Marshal(routine)
			req, _ := http.NewRequest(http.MethodPost, routinesURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", bearerToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Get free schedules
		url := fmt.Sprintf("%s?users=%s", freeSchedulesURL, userId)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		actual := unmarshallFreeSchedule(t, w.Body.Bytes())
		assert.NotEmpty(t, actual.Schedules, "Should return free schedules")

		// Verify Monday 8-10 and Wednesday 18-19 are not in free schedules
		for _, schedule := range actual.Schedules {
			if schedule.Day == "Monday" {
				// Check that 8-10 hours are not included in free schedules
				assert.False(t, (schedule.StartHour <= 8 && schedule.EndHour > 8) ||
					(schedule.StartHour >= 8 && schedule.StartHour < 10),
					"Hours 8-10 on Monday should not be free")
			}
			if schedule.Day == "Wednesday" {
				// Check that 18-19 hours are not included in free schedules
				assert.False(t, (schedule.StartHour <= 18 && schedule.EndHour > 18),
					"Hour 18-19 on Wednesday should not be free")
			}
		}
	})

	t.Run("GET /users/freeSchedules/ - Multiple Users", func(t *testing.T) {
		defer test.ClearAllData()

		// Create some routines
		routine := model.RoutineDTO{
			UserID:      userId,
			Name:        "Morning Workout",
			Description: "Cardio and strength training",
			Schedule: model.Schedule{
				Day:       "Monday",
				StartHour: 8,
				EndHour:   10,
			},
		}

		body, _ := json.Marshal(routine)
		req, _ := http.NewRequest(http.MethodPost, routinesURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Get free schedules for single user (which we have)
		url := fmt.Sprintf("%s?users=%s", freeSchedulesURL, userId)
		req, _ = http.NewRequest(http.MethodGet, url, nil)
		req.Header.Add("Authorization", bearerToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		response := unmarshallFreeSchedule(t, w.Body.Bytes())
		assert.NotEmpty(t, response.Schedules, "Should return free schedules for single user")
		assert.Len(t, response.Schedules, 8, "Should return eight free schedule for the user")
	})

	t.Run("GET /users/freeSchedules/ - No Users Provided", func(t *testing.T) {
		defer test.ClearAllData()

		req, _ := http.NewRequest(http.MethodGet, freeSchedulesURL, nil)
		req.Header.Add("Authorization", bearerToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

		response := unmarshallFreeSchedule(t, w.Body.Bytes())
		assert.Empty(t, response.Schedules, "Should return empty schedules when no users provided")
	})

	t.Run("GET /users/freeSchedules/ - Unauthorized", func(t *testing.T) {
		defer test.ClearAllData()

		url := fmt.Sprintf("%s?users=%s", freeSchedulesURL, userId)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		// No authorization token
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

		expected := model.ErrorRfc9457{
			Title:    "Unauthorized user",
			Detail:   `The user isn't authorized because no Authorization header is provided`,
			Status:   http.StatusUnauthorized,
			Type:     "about:blank",
			Instance: url,
		}

		var response model.ErrorRfc9457
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, expected, response)
	})
}
