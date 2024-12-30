package handler

import (
	"github.com/labstack/echo/v4"
	"go-echo/repository"
	"net/http"
	"strconv"
)

type OrganizationHandlerImpl struct {
	OrganizationRepository repository.OrganizationRepository
}

type OrganizationHandler interface {
	CreateOrganization(c echo.Context) error
	ReadOrganization(c echo.Context) error
	EditOrganization(c echo.Context) error
	DeleteOrganization(c echo.Context) error
	AllOrganization(c echo.Context) error
}

type Params struct {
	Id int `param:"id"`
}

type CreateRequestBody struct {
	Name     string `json:"name" validate:"required"`
	ParentId *int   `json:"parent_id"`
}

type ResponseSuccess struct {
	Data *OrganizationDataEntity `json:"data"`
}
type ResponseFailed struct {
	Message string `json:"message"`
}

type OrganizationDataEntity struct {
	Id       int                       `json:"id"`
	Name     string                    `json:"name"`
	ParentId *int                      `json:"parent_id"`
	Children []*OrganizationDataEntity `json:"children"`
}

type Data struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId *int   `json:"parent_id"`
}

func (s *OrganizationHandlerImpl) CreateOrganization(c echo.Context) error {
	req := new(CreateRequestBody)
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFailed{Message: err.Error()})
	}

	level := 0
	if req.ParentId != nil {
		org, err := s.OrganizationRepository.Get(c.Request().Context(), *req.ParentId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
		}
		if org.Level == 4 {
			return c.JSON(http.StatusBadRequest, ResponseFailed{Message: "Maximum nodes 5 level"})
		}
		level = org.Level + 1
	}

	out, err := s.OrganizationRepository.Create(c.Request().Context(), repository.OrganizationEntity{
		Name:     req.Name,
		ParentId: req.ParentId,
		Level:    level,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	resp := &OrganizationDataEntity{
		Id:       out.Id,
		Name:     out.Name,
		ParentId: out.ParentId,
	}

	return c.JSON(http.StatusOK, ResponseSuccess{Data: resp})
}

func (s *OrganizationHandlerImpl) ReadOrganization(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFailed{Message: err.Error()})
	}

	out, err := s.OrganizationRepository.Get(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	var resp OrganizationDataEntity
	resp.Id = out.Id
	resp.Name = out.Name
	resp.ParentId = out.ParentId

	outs, err := s.OrganizationRepository.GetByParentArr(c.Request().Context(), []int{id})
	for _, v := range outs {
		resp.Children = append(resp.Children, &OrganizationDataEntity{
			Id:       v.Id,
			Name:     v.Name,
			ParentId: v.ParentId,
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSuccess{Data: &resp})
}

func (s *OrganizationHandlerImpl) EditOrganization(c echo.Context) error {
	req := new(CreateRequestBody)
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFailed{Message: err.Error()})
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFailed{Message: err.Error()})
	}

	out, err := s.OrganizationRepository.Update(c.Request().Context(), repository.OrganizationEntity{
		Name:     req.Name,
		ParentId: req.ParentId,
		Id:       id,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	resp := &OrganizationDataEntity{
		Id:       out.Id,
		Name:     out.Name,
		ParentId: out.ParentId,
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Data: resp})
}

func (s *OrganizationHandlerImpl) DeleteOrganization(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseFailed{Message: err.Error()})
	}

	ids := []int{id}
	for i := 0; i < 4; i++ {
		outs, err := s.OrganizationRepository.GetByParentArr(c.Request().Context(), ids)
		for _, v := range outs {
			ids = append(ids, v.Id)
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
		}
	}

	err = s.OrganizationRepository.Delete(c.Request().Context(), ids)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSuccess{})
}

func (s *OrganizationHandlerImpl) AllOrganization(c echo.Context) error {

	outs, err := s.OrganizationRepository.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseFailed{Message: err.Error()})
	}

	var resp []Data

	for _, v := range outs {
		resp = append(resp, Data{
			Id:       v.Id,
			Name:     v.Name,
			ParentId: v.ParentId,
		})
	}

	return c.JSON(http.StatusOK, &resp)
}

type OptsParams struct {
	OrganizationRepository repository.OrganizationRepository
}

func NewOrganizations(opts OptsParams) OrganizationHandler {
	return &OrganizationHandlerImpl{
		OrganizationRepository: opts.OrganizationRepository,
	}
}
