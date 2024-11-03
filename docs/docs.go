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
        "/api/v1/weather": {
            "get": {
                "description": "Get weather data for a specific city.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "weather"
                ],
                "summary": "Get weather by city",
                "parameters": [
                    {
                        "type": "string",
                        "description": "City name",
                        "name": "city",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Weather"
                        }
                    },
                    "400": {
                        "description": "Validation error",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    },
                    "404": {
                        "description": "City not found",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    },
                    "504": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    }
                }
            }
        },
        "/api/v1/weather/feedback": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Submit feedback about the weather in a specific city.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "weather"
                ],
                "summary": "Submit weather feedback",
                "parameters": [
                    {
                        "description": "Weather feedback",
                        "name": "feedback",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.WeatherFeedbackReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.WeatherFeedbackResp"
                        }
                    },
                    "400": {
                        "description": "Invalid request data",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/result.Err"
                        }
                    }
                }
            }
        },
        "/healthz": {
            "get": {
                "description": "This endpoint returns the health status of the application.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ready": {
            "get": {
                "description": "This endpoint returns the readiness status of the application.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Readiness check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.WeatherFeedbackReq": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "dto.WeatherFeedbackResp": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.Astro": {
            "type": "object",
            "properties": {
                "sunrise": {
                    "type": "string"
                },
                "sunset": {
                    "type": "string"
                }
            }
        },
        "model.Current": {
            "type": "object",
            "properties": {
                "cloud": {
                    "type": "integer"
                },
                "condition": {
                    "type": "string"
                },
                "feelslike_c": {
                    "type": "number"
                },
                "heatindex_c": {
                    "type": "number"
                },
                "humidity": {
                    "type": "integer"
                },
                "last_updated": {
                    "type": "string"
                },
                "precip_mm": {
                    "type": "integer"
                },
                "pressure_mb": {
                    "type": "integer"
                },
                "temp_c": {
                    "type": "number"
                },
                "uv": {
                    "type": "integer"
                },
                "vis_km": {
                    "type": "integer"
                },
                "wind_degree": {
                    "type": "integer"
                },
                "wind_dir": {
                    "type": "string"
                },
                "wind_kph": {
                    "type": "number"
                }
            }
        },
        "model.Location": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "lat": {
                    "type": "number"
                },
                "localtime": {
                    "type": "string"
                },
                "lon": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "tz_id": {
                    "type": "string"
                },
                "tz_offset": {
                    "type": "integer"
                }
            }
        },
        "model.Weather": {
            "type": "object",
            "properties": {
                "astro": {
                    "$ref": "#/definitions/model.Astro"
                },
                "current": {
                    "$ref": "#/definitions/model.Current"
                },
                "location": {
                    "$ref": "#/definitions/model.Location"
                }
            }
        },
        "result.Err": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "detail": {
                    "type": "string"
                },
                "title": {
                    "$ref": "#/definitions/result.ErrTitle"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "result.ErrTitle": {
            "type": "string",
            "enum": [
                "Validation problem",
                "Not Found",
                "Conflict",
                "Unauthorized",
                "Request Timeout"
            ],
            "x-enum-varnames": [
                "Validation",
                "NotFound",
                "Conflict",
                "UnAuthorized",
                "GatewayTimeout"
            ]
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
