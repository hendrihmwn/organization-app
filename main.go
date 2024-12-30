package main

import (
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go-echo/handler"
	"go-echo/repository"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())

	db, err := sql.Open("postgres", "postgres://postgres:12345678@localhost:5432/organizations?sslmode=disable")
	if err != nil {
		panic(err)
	}

	repositoryOrganization := repository.NewOrganizationRepository(db)

	handlerOrganization := handler.NewOrganizations(handler.OptsParams{OrganizationRepository: repositoryOrganization})
	e.POST("/organization", handlerOrganization.CreateOrganization)
	e.PUT("/organization/:id", handlerOrganization.EditOrganization)
	e.GET("/organization/:id", handlerOrganization.ReadOrganization)
	e.DELETE("/organization/:id", handlerOrganization.DeleteOrganization)
	e.GET("/organizations", handlerOrganization.AllOrganization)

	e.Logger.Fatal(e.Start(":3000"))
}
