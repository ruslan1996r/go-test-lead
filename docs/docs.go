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
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/clients": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "Receives a list of clients",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/storage.Client"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "Creates a new client",
                "parameters": [
                    {
                        "description": "New client payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.ClientRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "bool"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/clients/assign": {
            "post": {
                "description": "Selects a suitable client for assignment. Assigns a Lead to him and returns ID of this client.\nInitially sort users by their availability and suitable time frames.\nThen sort users by their priority and percentage of free capacity. Select the user with the highest indicator.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "Assigns a Lead to a suitable client",
                "parameters": [
                    {
                        "description": "Assign lead payload",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storage.AssignLeadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/storage.Lead"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/clients/{id}": {
            "get": {
                "description": "Returns a single client array with the found user or a string with an error in case user is not found",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "client"
                ],
                "summary": "Get client by clientID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Client ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/storage.Client"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "storage.AssignLeadRequest": {
            "type": "object",
            "properties": {
                "lead_end": {
                    "type": "string"
                },
                "lead_start": {
                    "type": "string"
                }
            }
        },
        "storage.Client": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lead_capacity": {
                    "type": "integer"
                },
                "leads": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.Lead"
                    }
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "type": "string",
                    "enum": [
                        "HIGH",
                        "MEDIUM",
                        "LOW"
                    ]
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "storage.ClientRequest": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "lead_capacity": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "priority": {
                    "type": "string",
                    "enum": [
                        "HIGH",
                        "MEDIUM",
                        "LOW"
                    ]
                },
                "start_date": {
                    "type": "string"
                }
            }
        },
        "storage.Lead": {
            "type": "object",
            "properties": {
                "client_id": {
                    "type": "integer"
                },
                "lead_end": {
                    "type": "string"
                },
                "lead_id": {
                    "type": "string"
                },
                "lead_start": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
