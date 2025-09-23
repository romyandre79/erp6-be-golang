package auth

import (
	"erp6-be-golang/internal/core/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterAuthModule(s *server.Server) {
	r := s.Engine.Group("/auth")
	{
		r.POST("/login", loginHandler(s))
		r.POST("/register", registerHandler(s))
	}
}

func loginHandler(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement JWT login
		c.JSON(http.StatusOK, gin.H{"message": "login success (dummy)"})
	}
}

func registerHandler(s *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement user registration
		c.JSON(http.StatusOK, gin.H{"message": "register success (dummy)"})
	}
}
