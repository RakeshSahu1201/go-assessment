package main

import (
	"log"
	"main/auth"
	"main/controllers"
	"main/db"
	"net/http"
	"os"

	_ "main/docs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	db.InitDB()
	defer db.CloseDB()

	// Get session secret and port from env
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("SESSION_SECRET is required in .env or system environment variables")
	}

	r := gin.Default()

	// Configure session middleware
	store := cookie.NewStore([]byte(sessionSecret))
	store.Options(sessions.Options{
		MaxAge:   3600, // 1-hour session timeout
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Set to true for HTTPS connections only
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("session", store))

	// Setup Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Login route
	r.POST("/login", auth.Login)
	r.POST("/logout", auth.Logout)

	// Receptionist routes
	receptionist := r.Group("/receptionist")
	receptionist.Use(auth.AuthMiddleware("receptionist"))
	receptionist.POST("/patients", controllers.AddPatientPOST)
	receptionist.GET("/patients", controllers.AllPatientsGET)
	receptionist.PUT("/patients/:id", controllers.UpdatePatientPUT)
	receptionist.DELETE("/patients/:id", controllers.PatientDELETE)

	// Doctor routes
	doctor := r.Group("/doctor")
	doctor.Use(auth.AuthMiddleware("doctor"))
	doctor.GET("/patients", controllers.AllPatientsForDoctorGET)
	doctor.PUT("/patients/:id/notes", controllers.UpdatePatientNotesPUT)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
