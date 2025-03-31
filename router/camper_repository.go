package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findCamperByIDHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	id := c.Param("id")

	camper, err := h.camperRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    camper,
	})
}

func (h *httpService) findAllCampersHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	var query model.CamperQueryInput

	if err := c.Bind(&query); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	campers, total, err := h.camperRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.Errorf("Error getting campers: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, withPaging(campers, total, query.PageOrDefault(), query.SizeOrDefault()))
}

func (h *httpService) updateCamperHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	id := c.Param("id")

	var camper model.Camper

	if err := c.Bind(&camper); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.IsCustomer() {
		logger.Errorf("User is not authorized to access this resource")
		return c.JSON(http.StatusForbidden, response{
			Success: false,
			Message: "forbidden",
		})
	}

	if err := h.camperRepo.Update(c.Request().Context(), id, camper); err != nil {
		logger.Errorf("Error updating camper: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    camper,
	})
}

func (h *httpService) createCamperHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	var camper model.Camper

	if err := c.Bind(&camper); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.IsCustomer() {
		logger.Errorf("User is not authorized to access this resource")
		return c.JSON(http.StatusForbidden, response{
			Success: false,
			Message: "forbidden",
		})
	}

	if err := h.camperRepo.Create(c.Request().Context(), camper); err != nil {
		logger.Errorf("Error creating camper: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, response{
		Success: true,
		Data:    camper,
	})
}

func (h *httpService) deleteCamperHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	id := c.Param("id")

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.IsCustomer() {
		logger.Errorf("User is not authorized to access this resource")
		return c.JSON(http.StatusForbidden, response{
			Success: false,
			Message: "forbidden",
		})
	}

	if err := h.camperRepo.Delete(c.Request().Context(), id); err != nil {
		logger.Errorf("Error deleting camper: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusNoContent, response{
		Success: true,
	})
}
