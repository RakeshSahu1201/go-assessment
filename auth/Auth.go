package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		role := session.Get("role")
		if role == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: no active session"})
			c.Abort()
			return
		}

		if requiredRole != "" && role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Login godoc
// @Summary Login to the system
// @Description This API endpoint allows users to log in by providing a username, password, and role.
// @Tags authentication
// @Param loginRequest body creds true "Login credentials"
// @Success 200 {object} string "Login successful"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Router /login [post]
func Login(c *gin.Context) {
	var loginRequest creds

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hardcoded authentication (Development Use Only)
	if (loginRequest.Username == "receptionist1" && loginRequest.Password == "reception123" && loginRequest.Role == "receptionist") ||
		(loginRequest.Username == "doctor1" && loginRequest.Password == "doctor123" && loginRequest.Role == "doctor") {

		// Create session
		session := sessions.Default(c)
		session.Set("username", loginRequest.Username)
		session.Set("role", loginRequest.Role)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

// Logout godoc
// @Summary Logout from the system
// @Description This API endpoint allows users to log out and clear their session.
// @Tags authentication
// @Success 200 {object} string "Logged out successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to logout"
// @Router /logout [post]
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear() // Clear all session data
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
