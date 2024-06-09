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
        "/app": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Get all key-value pairs",
                "operationId": "get-all",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetAllResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
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
                    "application"
                ],
                "summary": "Set value to key",
                "operationId": "set-value-to-key",
                "parameters": [
                    {
                        "description": "Key and value to set",
                        "name": "key",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SetRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Delete all key-value pairs",
                "operationId": "delete-all",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            },
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Append value to key",
                "operationId": "append-value-to-key",
                "parameters": [
                    {
                        "description": "Key and value to append",
                        "name": "key",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.AppendRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/app/{key}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Get value by key",
                "operationId": "get-value-by-key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key to get",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Delete value by key",
                "operationId": "delete-value-by-key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key to delete",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DeleteResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/app/{key}/strlen": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "application"
                ],
                "summary": "Get value string length by key",
                "operationId": "get-value-string-length-by-key",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Key to get",
                        "name": "key",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.StrlenResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/cluster/log": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cluster"
                ],
                "summary": "Request cluster leader's log",
                "operationId": "request-log",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.RequestLogResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/cluster/ping": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cluster"
                ],
                "summary": "Ping cluster",
                "operationId": "ping-cluster",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.PingResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.ResponseMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "data.LogEntry": {
            "type": "object",
            "properties": {
                "command": {
                    "type": "string"
                },
                "term": {
                    "type": "integer"
                },
                "value": {
                    "description": "Optional field (maybe zero Value)",
                    "type": "string"
                }
            }
        },
        "handler.AppendRequest": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "handler.DeleteResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/utils.KeyVal"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.GetAllResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/utils.KeyVal"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.GetResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/utils.KeyVal"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.PingResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.RequestLogResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.LogEntry"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "handler.SetRequest": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "handler.StrlenResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "utils.KeyVal": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "utils.ResponseMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Web Client API",
	Description:      "This is the API documentation for the web client of the application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
