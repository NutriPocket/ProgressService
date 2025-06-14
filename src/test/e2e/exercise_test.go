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

// Helper functions to unmarshall responses
func unmarshallExerciseData(t *testing.T, body []byte) model.ExerciseData {
    var response map[string]any
    err := json.Unmarshal(body, &response)
    assert.NoError(t, err, "Response should be valid JSON")

    bytes, err := json.Marshal(response["data"])
    assert.NoError(t, err, "Response.data should be valid JSON")

    var actual model.ExerciseData
    err = json.Unmarshal(bytes, &actual)
    assert.NoError(t, err, "Response.data should be valid JSON")

    return actual
}

func unmarshallAllExercisesInDay(t *testing.T, body []byte) model.AllExercisesInDay {
    var response map[string]any
    err := json.Unmarshal(body, &response)
    assert.NoError(t, err, "Response should be valid JSON")

    bytes, err := json.Marshal(response["data"])
    assert.NoError(t, err, "Response.data should be valid JSON")

    var actual model.AllExercisesInDay
    err = json.Unmarshal(bytes, &actual)
    assert.NoError(t, err, "Response.data should be valid JSON")

    return actual
}

// Test POST /users/:userId/exercises/
func TestCreateExercise(t *testing.T) {
    userId := testUser.ID
    baseURL := fmt.Sprintf("/users/%s/exercises/", userId)

    t.Run("POST /users/:userId/exercises/ - Create Exercise Successfully", func(t *testing.T) {
        defer test.ClearAllData()

        payload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
        }

        body, _ := json.Marshal(payload)
        req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201")

        actual := unmarshallExerciseData(t, w.Body.Bytes())

        assert.Equal(t, payload.UserID, actual.UserID)
        assert.Equal(t, payload.ExerciseName, actual.ExerciseName)
        assert.Equal(t, payload.CaloriesBurned, actual.CaloriesBurned)
        assert.NotEmpty(t, actual.ID)
        assert.NotEmpty(t, actual.CreatedAt)
    })

    t.Run("POST /users/:userId/exercises/ - Missing Required Fields", func(t *testing.T) {
        defer test.ClearAllData()

        // Missing required field (exerciseName)
        payload := map[string]interface{}{
            "userId":         userId,
            "caloriesBurned": 450.5,
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
        assert.Equal(t, "Invalid exercise data", response.Title)
    })

    t.Run("POST /users/:userId/exercises/ - Invalid Calories Value", func(t *testing.T) {
        defer test.ClearAllData()

        // Negative calories value
        payload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: -50,
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
        assert.Equal(t, "Invalid exercise data", response.Title)
    })

    t.Run("POST /users/:userId/exercises/ - Unauthorized", func(t *testing.T) {
        defer test.ClearAllData()

        payload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
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
}

// Test GET /users/:userId/exercises/
func TestGetExercises(t *testing.T) {
    userId := testUser.ID
    baseURL := fmt.Sprintf("/users/%s/exercises/", userId)

    t.Run("GET /users/:userId/exercises/ - Get Today's Exercises", func(t *testing.T) {
        defer test.ClearAllData()

        // Create some exercises first
        exercises := []model.ExerciseDTO{
            {
                UserID:         userId,
                ExerciseName:   "Running",
                CaloriesBurned: 450.5,
            },
            {
                UserID:         userId,
                ExerciseName:   "Weightlifting",
                CaloriesBurned: 300.0,
            },
        }

        for _, exercise := range exercises {
            body, _ := json.Marshal(exercise)
            req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Add("Authorization", bearerToken)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)
        }

        // Get today's exercises (default behavior)
        req, _ := http.NewRequest(http.MethodGet, baseURL, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

        actual := unmarshallAllExercisesInDay(t, w.Body.Bytes())

        assert.Len(t, actual.Exercises, 2, "Should return 2 exercises")
        assert.Equal(t, 750.5, actual.TotalBurned, "Total burned calories should be correct")

        // Verify exercise data
        foundExercises := 0
        for _, e := range actual.Exercises {
            if e.ExerciseName == "Running" {
                foundExercises++
                assert.Equal(t, 450.5, e.CaloriesBurned)
            } else if e.ExerciseName == "Weightlifting" {
                foundExercises++
                assert.Equal(t, 300.0, e.CaloriesBurned)
            }
        }
        assert.Equal(t, 2, foundExercises, "Should find both created exercises")
    })

    t.Run("GET /users/:userId/exercises/?date=YYYY-MM-DD - Get Exercises for Specific Date", func(t *testing.T) {
        defer test.ClearAllData()

        // Create an exercise
        exercise := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
        }

        body, _ := json.Marshal(exercise)
        req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        // Get exercises for today's date explicitly
        today := time.Now().Format("2006-01-02")
        url := fmt.Sprintf("%s?date=%s", baseURL, today)
        req, _ = http.NewRequest(http.MethodGet, url, nil)
        req.Header.Add("Authorization", bearerToken)
        w = httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

        actual := unmarshallAllExercisesInDay(t, w.Body.Bytes())

        assert.Len(t, actual.Exercises, 1, "Should return 1 exercise")
        assert.Equal(t, 450.5, actual.TotalBurned, "Total burned calories should be correct")
        assert.Equal(t, "Running", actual.Exercises[0].ExerciseName)
    })

    t.Run("GET /users/:userId/exercises/?date=YYYY-MM-DD - No Exercises for Date", func(t *testing.T) {
        defer test.ClearAllData()

        // Get exercises for a date with no data
        url := fmt.Sprintf("%s?date=2020-01-01", baseURL)
        req, _ := http.NewRequest(http.MethodGet, url, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

        actual := unmarshallAllExercisesInDay(t, w.Body.Bytes())

        assert.Empty(t, actual.Exercises, "Should return empty exercises array")
        assert.Equal(t, 0.0, actual.TotalBurned, "Total burned calories should be zero")
    })

    t.Run("GET /users/:userId/exercises/?date=invalid - Invalid Date Format", func(t *testing.T) {
        defer test.ClearAllData()

        // Get exercises with invalid date format
        url := fmt.Sprintf("%s?date=invalid-date", baseURL)
        req, _ := http.NewRequest(http.MethodGet, url, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, "Invalid date format", response.Title)
    })

    t.Run("GET /users/:userId/exercises/ - Unauthorized", func(t *testing.T) {
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

// Test PUT /users/:userId/exercises/:id
func TestUpdateExercise(t *testing.T) {
    userId := testUser.ID
    baseURL := fmt.Sprintf("/users/%s/exercises/", userId)

    t.Run("PUT /users/:userId/exercises/:id - Update Exercise Successfully", func(t *testing.T) {
        defer test.ClearAllData()

        // Create an exercise first
        createPayload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
        }

        var exerciseId uint64
        {
            body, _ := json.Marshal(createPayload)
            req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Add("Authorization", bearerToken)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            createdExercise := unmarshallExerciseData(t, w.Body.Bytes())
            exerciseId = createdExercise.ID
        }

        // Update the exercise
        updatePayload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Sprint Running",
            CaloriesBurned: 500.0,
        }

        updateURL := fmt.Sprintf("%s%d", baseURL, exerciseId)
        body, _ := json.Marshal(updatePayload)
        req, _ := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

        actual := unmarshallExerciseData(t, w.Body.Bytes())

        assert.Equal(t, exerciseId, actual.ID)
        assert.Equal(t, updatePayload.UserID, actual.UserID)
        assert.Equal(t, updatePayload.ExerciseName, actual.ExerciseName)
        assert.Equal(t, updatePayload.CaloriesBurned, actual.CaloriesBurned)
    })

    t.Run("PUT /users/:userId/exercises/:id - Exercise Not Found", func(t *testing.T) {
        defer test.ClearAllData()

        // Try to update a non-existent exercise
        updatePayload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Sprint Running",
            CaloriesBurned: 500.0,
        }

        updateURL := fmt.Sprintf("%s999999", baseURL) // Non-existent ID
        body, _ := json.Marshal(updatePayload)
        req, _ := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, "Exercise not found", response.Title)
    })

    t.Run("PUT /users/:userId/exercises/:id - Invalid Exercise Data", func(t *testing.T) {
        defer test.ClearAllData()

        // Create an exercise first
        createPayload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
        }

        var exerciseId uint64
        {
            body, _ := json.Marshal(createPayload)
            req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Add("Authorization", bearerToken)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            createdExercise := unmarshallExerciseData(t, w.Body.Bytes())
            exerciseId = createdExercise.ID
        }

        // Try to update with invalid data (missing required field)
        updatePayload := map[string]interface{}{
            "userId":         userId,
            "caloriesBurned": 500.0,
            // Missing exerciseName
        }

        updateURL := fmt.Sprintf("%s%d", baseURL, exerciseId)
        body, _ := json.Marshal(updatePayload)
        req, _ := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, "Invalid exercise data", response.Title)
    })

    t.Run("PUT /users/:userId/exercises/:id - Unauthorized", func(t *testing.T) {
        defer test.ClearAllData()

        updateURL := fmt.Sprintf("%s123", baseURL)
        req, _ := http.NewRequest(http.MethodPut, updateURL, nil)
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
            Instance: updateURL,
        }

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, expected, response)
    })
}

// Test DELETE /users/:userId/exercises/:id
func TestDeleteExercise(t *testing.T) {
    userId := testUser.ID
    baseURL := fmt.Sprintf("/users/%s/exercises/", userId)

    t.Run("DELETE /users/:userId/exercises/:id - Delete Exercise Successfully", func(t *testing.T) {
        defer test.ClearAllData()

        // Create an exercise first
        createPayload := model.ExerciseDTO{
            UserID:         userId,
            ExerciseName:   "Running",
            CaloriesBurned: 450.5,
        }

        var exerciseId uint64
        {
            body, _ := json.Marshal(createPayload)
            req, _ := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            req.Header.Add("Authorization", bearerToken)
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            createdExercise := unmarshallExerciseData(t, w.Body.Bytes())
            exerciseId = createdExercise.ID
        }

        // Delete the exercise
        deleteURL := fmt.Sprintf("%s%d", baseURL, exerciseId)
        req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNoContent, w.Code, "Status code should be 204")

        // Verify that the exercise is deleted by trying to get it
        today := time.Now().Format("2006-01-02")
        getURL := fmt.Sprintf("%s?date=%s", baseURL, today)
        req, _ = http.NewRequest(http.MethodGet, getURL, nil)
        req.Header.Add("Authorization", bearerToken)
        w = httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200")

        actual := unmarshallAllExercisesInDay(t, w.Body.Bytes())
        assert.Empty(t, actual.Exercises, "Should return empty exercises array after deletion")
        assert.Equal(t, 0.0, actual.TotalBurned, "Total burned calories should be zero")
    })

    t.Run("DELETE /users/:userId/exercises/:id - Exercise Not Found", func(t *testing.T) {
        defer test.ClearAllData()

        // Try to delete a non-existent exercise
        deleteURL := fmt.Sprintf("%s999999", baseURL) // Non-existent ID
        req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404")

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, "Exercise not found", response.Title)
    })

    t.Run("DELETE /users/:userId/exercises/:id - Invalid ID Format", func(t *testing.T) {
        defer test.ClearAllData()

        // Try to delete with invalid ID format
        deleteURL := fmt.Sprintf("%sinvalid", baseURL) // Non-numeric ID
        req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
        req.Header.Add("Authorization", bearerToken)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400")

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, "Invalid exercise ID", response.Title)
    })

    t.Run("DELETE /users/:userId/exercises/:id - Unauthorized", func(t *testing.T) {
        defer test.ClearAllData()

        deleteURL := fmt.Sprintf("%s123", baseURL)
        req, _ := http.NewRequest(http.MethodDelete, deleteURL, nil)
        // No authorization token
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusUnauthorized, w.Code, "Status code should be 401")

        expected := model.ErrorRfc9457{
            Title:    "Unauthorized user",
            Detail:   `The user isn't authorized because no Authorization header is provided`,
            Status:   http.StatusUnauthorized,
            Type:     "about:blank",
            Instance: deleteURL,
        }

        var response model.ErrorRfc9457
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err, "Response should be valid JSON")
        assert.Equal(t, expected, response)
    })
}