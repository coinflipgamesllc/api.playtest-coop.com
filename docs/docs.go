// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "https://playtest-coop.com/terms-of-service",
        "contact": {
            "name": "Coin Flip Games",
            "email": "hi@coinflipgames.co"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authenticate a user",
                "parameters": [
                    {
                        "description": "User email/password combo",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.UserTokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/signup": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Create and authenticates a new user",
                "parameters": [
                    {
                        "description": "User name, email, and password",
                        "name": "credentials",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.SignupRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/controller.UserTokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/token": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Regenerate the access token and refresh token, given a valid refresh token.",
                "parameters": [
                    {
                        "description": "Refresh token originally acquired from /auth/token, /auth/signup, or /auth/login",
                        "name": "refresh_token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    }
                }
            }
        },
        "/auth/user": {
            "get": {
                "description": "The authentication token includes the user's ID as the subject. We extract that and use it to pull the user from the database.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Retrieve the authenticated user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetUserResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.UnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Update authenticated user",
                "parameters": [
                    {
                        "description": "User data to update",
                        "name": "params",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/controller.UpdateUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.GetUserResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.UnauthorizedResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/files": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "List files belonging to the authenticated user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.ListUserFilesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Save a record of a file stored in S3",
                "parameters": [
                    {
                        "description": "File data",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateFileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.AckResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/files/:id": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "remove a file by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "File ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.AckResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            }
        },
        "/files/sign": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "Generate a presigned URL for the client to upload directly to S3",
                "parameters": [
                    {
                        "description": "File data",
                        "name": "file",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.PresignUploadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controller.PresignUploadResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.RequestErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.ServerErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.AckResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "controller.CreateFileRequest": {
            "type": "object",
            "required": [
                "filename",
                "object",
                "role",
                "size"
            ],
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "What a cool image of a game!"
                },
                "filename": {
                    "type": "string",
                    "example": "example-image.png"
                },
                "game": {
                    "type": "integer",
                    "example": 123
                },
                "object": {
                    "type": "string",
                    "example": "asd9fhgaoseucgewio.png"
                },
                "role": {
                    "type": "string",
                    "example": "Image"
                },
                "size": {
                    "type": "integer",
                    "example": 1241231
                }
            }
        },
        "controller.GetUserResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/domain.User"
                }
            }
        },
        "controller.ListUserFilesResponse": {
            "type": "object",
            "properties": {
                "files": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.File"
                    }
                }
            }
        },
        "controller.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "password": {
                    "type": "string",
                    "example": "AVerySecurePassword123!"
                }
            }
        },
        "controller.PresignUploadRequest": {
            "type": "object",
            "required": [
                "extension",
                "name"
            ],
            "properties": {
                "extension": {
                    "type": "string",
                    "example": "jpg"
                },
                "name": {
                    "type": "string",
                    "example": "my-awesome-file.jpg"
                }
            }
        },
        "controller.PresignUploadResponse": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string",
                    "example": "https://assets.playtest-coop.com/..."
                }
            }
        },
        "controller.RefreshTokenRequest": {
            "type": "object",
            "required": [
                "refresh_token"
            ],
            "properties": {
                "refresh_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"
                }
            }
        },
        "controller.RequestErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controller.ServerErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controller.SignupRequest": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "User McUserton"
                },
                "password": {
                    "type": "string",
                    "example": "AVerySecurePassword123!"
                }
            }
        },
        "controller.TokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"
                },
                "refresh_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"
                }
            }
        },
        "controller.UnauthorizedResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "controller.UpdateUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "user@example.com"
                },
                "name": {
                    "type": "string",
                    "example": "User McUserton"
                },
                "new_password": {
                    "type": "string",
                    "example": "AVerySecurePassword123!"
                },
                "old_password": {
                    "type": "string",
                    "example": "NotASecurePassword"
                },
                "pronouns": {
                    "type": "string",
                    "example": "they/them"
                }
            }
        },
        "controller.UserTokenResponse": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgwNTY5NzksIm5hbWUiOiJSb2IgTmV3dG9uIiwic3ViIjoxfQ.KKUtLne51DqBPqQxZZmCFsjsGAeYRukZNcXCx6IpLN8"
                },
                "refresh_token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDgxNDA1MTcsInN1YiI6MX0.D5kR_AxkqIN6xCxvP07ZUIfYxbfdTrXAe7J03nGvkPw"
                },
                "user": {
                    "$ref": "#/definitions/domain.User"
                }
            }
        },
        "domain.File": {
            "type": "object",
            "properties": {
                "caption": {
                    "type": "string",
                    "example": "What a cool image of a game!"
                },
                "created_at": {
                    "type": "string",
                    "example": "2020-12-11T15:29:49.321629-08:00"
                },
                "filename": {
                    "type": "string",
                    "example": "example-image.png"
                },
                "id": {
                    "type": "integer",
                    "example": 123
                },
                "role": {
                    "type": "string",
                    "example": "Image"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2020-12-13T15:42:40.578904-08:00"
                },
                "url": {
                    "type": "string",
                    "example": "https://assets.playtest-coop.com/asd9fhgaoseucgewio.png"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2020-12-11T15:29:49.321629-08:00"
                },
                "id": {
                    "type": "integer",
                    "example": 123
                },
                "name": {
                    "type": "string",
                    "example": "User McUserton"
                },
                "pronouns": {
                    "type": "string",
                    "example": "they/them"
                },
                "updated_at": {
                    "type": "string",
                    "example": "2020-12-13T15:42:40.578904-08:00"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "api.playtest-coop.com",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "Playtest Co-op API",
	Description: "This is the backend for all Playtest Co-op related data",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
