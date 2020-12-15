package controller

import "github.com/gin-gonic/gin"

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
