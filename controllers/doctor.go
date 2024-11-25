package controllers

import (
	"main/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AllPatientsForDoctorGET godoc
// @Summary Get all patients for a doctor
// @Description This API endpoint returns a list of all patients assigned to a doctor.
// @Tags doctor
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PatientSwagger "List of patients"
// @Failure 500 {object} models.ErrorResponse
// @Router /doctor/patients [get]
func AllPatientsForDoctorGET(c *gin.Context) {
	client := db.GetClient()
	patients, err := client.Patient.Query().All(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

type updatePatientNotesURIPUT struct {
	ID string `uri:"id" binding:"required"`
}

type updatePatientNotesRequestPUT struct {
	Notes string `json:"notes" binding:"required"`
}

// UpdatePatientNotesPUT godoc
// @Summary Update patient notes
// @Description This API endpoint updates the notes for a specific patient identified by their ID.
// @Tags doctor
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param request body updatePatientNotesRequestPUT true "Updated Notes"
// @Success 200 {object} models.PatientSwagger "Updated patient details"
// @Failure 400 {object} models.ErrorResponse "Invalid request (bad ID or JSON)"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /doctor/patients/{id}/notes [put]
// @Security BearerAuth
func UpdatePatientNotesPUT(c *gin.Context) {
	var patientURI updatePatientNotesURIPUT
	if err := c.ShouldBindUri(&patientURI); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID in URI"})
		return
	}

	patientID, err := strconv.Atoi(patientURI.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID must be a valid number"})
		return
	}

	var patientData updatePatientNotesRequestPUT
	if err := c.ShouldBindJSON(&patientData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid json data"})
		return
	}

	client := db.GetClient()
	patient, err := client.Patient.UpdateOneID(patientID).
		SetNotes(patientData.Notes).
		Save(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient notes"})
		return
	}

	c.JSON(http.StatusOK, patient)
}
