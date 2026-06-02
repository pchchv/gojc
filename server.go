package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// handlePing checks that the server is up and running.
func handlePing(c *echo.Context) error {
	message := "JSON collector API. Version 0.0.1"
	return c.String(http.StatusOK, message)
}

// handleCollect collects any dynamic data types
func handleCollect(c *echo.Context) error {
	rawParam := c.QueryParam("struct")
	if rawParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "parameter 'struct' is required"})
	}

	// Parsing the incoming query string
	parts := strings.Split(rawParam, ",")
	if len(parts) < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid format"})
	}

	structName := exportedName(parts[0])
	fieldsData := make(map[string]any)
	// Parse raw strings and automatically determine their deep type
	for _, part := range parts[1:] {
		pair := strings.SplitN(part, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			valStr := pair[1]
			fieldsData[key] = parseAnyType(valStr)
		}
	}

	// Dynamically generate field descriptions (StructField) for reflect
	var structFields []reflect.StructField
	for key, val := range fieldsData {
		fieldName := exportedName(key)
		fieldType := reflect.TypeOf(val)
		// If the type is not defined (nil),
		// make it a string by default
		if fieldType == nil {
			fieldType = reflect.TypeOf("")
		}

		structFields = append(structFields, reflect.StructField{
			Name: fieldName,
			Type: fieldType,
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s"`, key)),
		})
	}

	// Create a structure type in memory at runtime
	dynamicType := reflect.StructOf(structFields)

	// Allocate memory for the structure and fill it with values
	dynamicStructValue := reflect.New(dynamicType).Elem()
	for key, val := range fieldsData {
		fieldName := exportedName(key)
		field := dynamicStructValue.FieldByName(fieldName)
		if field.IsValid() && field.CanSet() {
			valOf := reflect.ValueOf(val)
			// Check if the type matches before assigning a value
			if field.Type() == valOf.Type() {
				field.Set(valOf)
			}
		}
	}

	// Extracting the finished structure
	finalStruct := dynamicStructValue.Interface()

	// Logs the name of the generated structure and its contents
	log.Printf("[Collector] The structure has been generated %s: %+v\n", structName, finalStruct)

	// Returns fully valid typed JSON to the client
	return c.JSON(http.StatusOK, finalStruct)
}

// The declaration of all routes comes from it.
func routes(e *echo.Echo) {
	e.GET("/", handlePing)
	e.GET("/ping", handlePing)
	e.GET("/collect", handleCollect)
}

func server() {
	e := echo.New()
	routes(e)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1000)))
	log.Fatal(e.Start(getEnvValue("HOST") + ":" + getEnvValue("PORT")).Error())
}
