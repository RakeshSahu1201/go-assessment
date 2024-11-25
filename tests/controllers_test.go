package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"main/auth"
	"main/controllers"
	"main/db"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddPatientPOST(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Ensure DB is initialized
	defer db.CloseDB() // Close DB after test

	// Initialize the session store
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))

	// Set up the routes
	router.POST("/login", auth.Login) // Login route to authenticate
	receptionist := router.Group("/receptionist")
	receptionist.Use(auth.AuthMiddleware("receptionist"))
	receptionist.POST("/patients", controllers.AddPatientPOST)

	loginBody := `{"username":"receptionist1", "password":"reception123", "role": "receptionist"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	// Verify login succeeded
	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	body := `{"name":"John Doe", "age":30, "gender":"male", "notes":"No notes"}`
	req, _ := http.NewRequest("POST", "/receptionist/patients", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie) // Attach the session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
}

func TestAllPatientsGET(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Ensure DB is initialized
	defer db.CloseDB() // Close DB after test

	setupTestData()
	// Initialize the session store
	store := cookie.NewStore([]byte("sessionSecret"))
	router.Use(sessions.Sessions("mysession", store))

	// Set up the routes
	router.POST("/login", auth.Login) // Login route to authenticate
	receptionist := router.Group("/receptionist")
	receptionist.Use(auth.AuthMiddleware("receptionist"))
	receptionist.GET("/patients", controllers.AllPatientsGET)

	loginBody := `{"username":"receptionist1", "password":"reception123", "role": "receptionist"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	// Verify login succeeded
	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	req, _ := http.NewRequest("GET", "/receptionist/patients", nil)
	req.Header.Set("Cookie", cookie) // Attach session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response) // Ensure the response has patient data
	assert.Equal(t, "John Doe", response[0]["name"])
}

func TestUpdatePatientPUT(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Ensure DB is initialized
	defer db.CloseDB() // Close DB after test

	id, _ := setupTestData()

	// Initialize the session store
	store := cookie.NewStore([]byte("sessionSecret"))
	router.Use(sessions.Sessions("mysession", store))

	// Set up the routes
	router.POST("/login", auth.Login) // Login route to authenticate
	receptionist := router.Group("/receptionist")
	receptionist.Use(auth.AuthMiddleware("receptionist"))
	receptionist.PUT("/patients/:id", controllers.UpdatePatientPUT)

	loginBody := `{"username":"receptionist1", "password":"reception123", "role": "receptionist"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	// Verify login succeeded
	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	body := `{"name":"John Doe Updated", "age":31, "gender":"male", "notes":"Updated notes"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/receptionist/patients/%d", id), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie) // Attach session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe Updated")
}

func TestDeletePatient(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Ensure DB is initialized
	defer db.CloseDB() // Close DB after test

	id, _ := setupTestData()
	// Initialize the session store
	store := cookie.NewStore([]byte("sessionSecret"))
	router.Use(sessions.Sessions("mysession", store))

	// Set up the routes
	router.POST("/login", auth.Login) // Login route to authenticate
	receptionist := router.Group("/receptionist")
	receptionist.Use(auth.AuthMiddleware("receptionist"))
	receptionist.DELETE("/patients/:id", controllers.PatientDELETE)

	loginBody := `{"username":"receptionist1", "password":"reception123", "role": "receptionist"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	// Verify login succeeded
	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/receptionist/patients/%d", id), nil)
	req.Header.Set("Cookie", cookie) // Attach session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Patient deleted")
}

func TestAllPatientsForDoctorGET(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Initialize the database
	defer db.CloseDB() // Close the database after the test

	setupTestData()
	// Initialize the session store
	store := cookie.NewStore([]byte("sessionSecret"))
	router.Use(sessions.Sessions("mysession", store))

	// Set up routes
	router.POST("/login", auth.Login)
	doctor := router.Group("/doctor")
	doctor.Use(auth.AuthMiddleware("doctor"))
	doctor.GET("/patients", controllers.AllPatientsForDoctorGET)

	loginBody := `{"username":"doctor1", "password":"doctor123", "role": "doctor"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	// Verify login succeeded
	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	req, _ := http.NewRequest("GET", "/doctor/patients", nil)
	req.Header.Set("Cookie", cookie) // Attach the session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response) // Ensure the response has patient data
	assert.Equal(t, "John Doe", response[0]["name"])
}

func TestUpdatePatientNotesPUT(t *testing.T) {
	// Create a new Gin router
	router := gin.Default()
	db.InitDBDefault() // Initialize the database
	defer db.CloseDB() // Close the database after the test

	id, _ := setupTestData()

	// Initialize the session store
	store := cookie.NewStore([]byte("sessionSecret"))
	router.Use(sessions.Sessions("mysession", store))

	// Set up routes
	router.POST("/login", auth.Login)
	doctor := router.Group("/doctor")
	doctor.Use(auth.AuthMiddleware("doctor"))
	doctor.PUT("/patients/:id/notes", controllers.UpdatePatientNotesPUT)

	loginBody := `{"username":"doctor1", "password":"doctor123", "role": "doctor"}`
	loginReq, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()

	router.ServeHTTP(loginResp, loginReq)

	assert.Equal(t, http.StatusOK, loginResp.Code)
	cookie := loginResp.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	body := `{"notes":"Updated patient notes"}`
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/doctor/patients/%d/notes", id), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie) // Attach the session cookie

	// Perform the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert the status code and response body
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated patient notes")
}

// helper function to create a patient
func setupTestData() (int, error) {
	// Create a new patient in the database
	client := db.GetClient()

	patient, err := client.Patient.Create().
		SetName("John Doe").
		SetAge(30).
		SetGender("Male").
		SetNotes("Initial notes").
		Save(context.Background())

	if err != nil {
		return 0, fmt.Errorf("failed to create patient: %v", err)
	}

	// Return the patient ID
	return patient.ID, nil
}
