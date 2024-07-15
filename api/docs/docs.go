// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "checks the user and returns tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "User Information to log in",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns access and refresh tokens",
                        "schema": {
                            "$ref": "#/definitions/auth.Tokens"
                        }
                    },
                    "401": {
                        "description": "if Access token fails it will returns this",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Something went wrong in server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/logout": {
            "post": {
                "description": "removes refresh token gets token from header",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "log outs user",
                "operationId": "logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "if Access token fails it will returns this",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Something went wrong in server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/refreshtoken": {
            "get": {
                "description": "generates new access token gets token from header",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "refreshes token",
                "operationId": "refresh",
                "responses": {
                    "200": {
                        "description": "if Access token fails it will returns this",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Something went wrong in server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "description": "Registers user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register User",
                "operationId": "register",
                "parameters": [
                    {
                        "description": "User information to create it",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "500": {
                        "description": "Something went wrong in server",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.Tokens": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "integer"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "auth.User": {
            "type": "object",
            "properties": {
                "bio": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "user_name": {
                    "type": "string"
                },
                "user_type": {
                    "type": "string"
                }
            }
        },
        "auth.UserLogin": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/auth",
	Schemes:          []string{},
	Title:            "ReserveDesk API",
	Description:      "Connecting artists and customers program",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}