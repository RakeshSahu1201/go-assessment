package controllers

import (
	"context"
	"main/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addPatientRequestPOST struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Gender string `json:"gender" binding:"required"`
	Notes  string `json:"notes"`
}

// AddPatientPOST godoc
// @Summary Add a new patient
// @Description This API endpoint creates a new patient record in the system.
// @Tags receptionist
// @Accept json
// @Produce json
// @Param request body addPatientRequestPOST true "Patient Details"
// @Success 201 {object} models.PatientSwagger "Details of the created patient"
// @Failure 400 {object} models.ErrorResponse "Invalid JSON request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /receptionist/patients [post]
// @Security BearerAuth
func AddPatientPOST(c *gin.Context) {
	var input addPatientRequestPOST
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json request"})
		return
	}

	client := db.GetClient()
	patient, err := client.Patient.Create().
		SetName(input.Name).
		SetAge(input.Age).
		SetGender(input.Gender).
		SetNotes(input.Notes).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create patient"})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// AllPatientsGET godoc
// @Summary Get all patients
// @Description This API endpoint retrieves all patients from the system.
// @Tags receptionist
// @Accept json
// @Produce json
// @Success 200 {array} models.PatientSwagger "List of all patients"
// @Failure 500 {object} models.ErrorResponse "Failed to fetch patients"
// @Router /receptionist/patients [get]
func AllPatientsGET(c *gin.Context) {
	client := db.GetClient()
	patients, err := client.Patient.Query().All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

type updatePatientRequestPUT struct {
	Name   string `json:"name" binding:"required"`
	Age    int    `json:"age" binding:"required"`
	Gender string `json:"gender" binding:"required"`
	Notes  string `json:"notes"`
}

type updatePatientURIPUT struct {
	ID string `uri:"id" binding:"required"`
}

// UpdatePatientPUT godoc
// @Summary Update a patient's details
// @Description This API endpoint updates the details of an existing patient using their ID.
// @Tags receptionist
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param patient body updatePatientRequestPUT true "Patient update data"
// @Success 200 {object} models.PatientSwagger "Patient details after update"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 500 {object} models.ErrorResponse "Failed to update patient"
// @Router /receptionist/patients/{id} [put]
func UpdatePatientPUT(c *gin.Context) {
	var uriParams updatePatientURIPUT
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID in URI"})
		return
	}

	patientID, err := strconv.Atoi(uriParams.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID must be a valid number"})
		return
	}

	var patientData updatePatientRequestPUT
	if err := c.ShouldBindJSON(&patientData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json request"})
		return
	}

	client := db.GetClient()

	patient, err := client.Patient.UpdateOneID(patientID).
		SetName(patientData.Name).
		SetAge(patientData.Age).
		SetGender(patientData.Gender).
		SetNotes(patientData.Notes).
		Save(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

type deletePatientURIRequestDELETE struct {
	ID string `uri:"id" binding:"required"`
}

// PatientDELETE godoc
// @Summary Delete a patient
// @Description This API endpoint deletes a patient by their ID.
// @Tags receptionist
// @Param id path string true "Patient ID"
// @Success 200 {object} string "Patient successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid patient ID"
// @Failure 500 {object} models.ErrorResponse "Failed to delete patient"
// @Router /receptionist/patients/{id} [delete]
func PatientDELETE(c *gin.Context) {
	var uriParams deletePatientURIRequestDELETE
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID in URI"})
		return
	}

	patientID, err := strconv.Atoi(uriParams.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID must be a valid number"})
		return
	}

	client := db.GetClient()

	if err := client.Patient.DeleteOneID(patientID).Exec(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete patient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted"})
}
