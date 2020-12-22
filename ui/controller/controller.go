package controller

import (
	"errors"

	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure/validation"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AckResponse simply acknowledges the request
type AckResponse struct {
	Message string `json:"message"`
}

func ackResponse(c *gin.Context) {
	c.JSON(200, AckResponse{Message: "ok"})
}

// NotFoundResponse to be paired with a 404
type NotFoundResponse struct {
	Error string `json:"error"`
}

func notFoundResponse(c *gin.Context, err string) {
	c.AbortWithStatusJSON(404, NotFoundResponse{Error: err})
}

// RequestErrorResponse to be paired with a 4xx
type RequestErrorResponse struct {
	Error string `json:"error"`
}

func requestErrorResponse(c *gin.Context, err string) {
	c.AbortWithStatusJSON(400, RequestErrorResponse{Error: err})
}

// ValidationErrorResponse for invalid requests
type ValidationErrorResponse struct {
	Errors map[string]string `json:"errors"`
}

func validationErrorResponse(c *gin.Context, err error) {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		f := validation.NewJSONFormatter()
		c.AbortWithStatusJSON(400, ValidationErrorResponse{Errors: f.Format(verr)})
		return
	}

	requestErrorResponse(c, err.Error())
}

// ServerErrorResponse to be paired with a 5xx
type ServerErrorResponse struct {
	Error string `json:"error"`
}

func serverErrorResponse(c *gin.Context, err string) {
	c.AbortWithStatusJSON(500, ServerErrorResponse{Error: err})
}

// UnauthorizedResponse to be paired with a 401/403
type UnauthorizedResponse struct {
	Error string `json:"error"`
}

func unauthorizedResponse(c *gin.Context) {
	c.AbortWithStatusJSON(401, UnauthorizedResponse{Error: "unauthorized"})
}

// userID helper function to extract the user's ID from the session
func userID(c *gin.Context) uint {
	// Retrieve the user ID from the session
	session := sessions.Default(c)
	id := session.Get("user_id")
	if id == nil {
		unauthorizedResponse(c)
		return 0
	}

	return id.(uint)
}
