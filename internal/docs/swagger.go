// Package docs provides Swagger API documentation for the Sutram platform.
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API server for the Sutram algorithmic trading platform",
        "title": "Sutram API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@sutram.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:9090",
    "basePath": "/api",
    "paths": {},
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header",
            "description": "Type \"Bearer\" followed by a space and the JWT token"
        }
    },
    "security": [
        {
            "BearerAuth": []
        }
    ],
    "tags": [
        {
            "description": "Strategy management endpoints",
            "name": "Strategy"
        },
        {
            "description": "User management endpoints",
            "name": "User"
        },
        {
            "description": "Instrument management endpoints",
            "name": "Instrument"
        },
        {
            "description": "Strategy deployment endpoints",
            "name": "Deployment"
        },
        {
            "description": "Strategy backtesting endpoints",
            "name": "Backtest"
        },
        {
            "description": "Marketplace endpoints",
            "name": "Marketplace"
        },
        {
            "description": "Community feature endpoints",
            "name": "Community"
        },
        {
            "description": "OAuth authentication endpoints",
            "name": "OAuth"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:9090",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "Sutram API",
	Description:      "API server for the Sutram algorithmic trading platform",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}