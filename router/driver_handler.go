package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findDriverByIDHandler(c echo.Context) error {
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

	driver, err := h.driverRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		logger.Errorf("Error querying driver: %v", err)
		return c.JSON(http.StatusNotFound, response{
			Success: false,
			Message: "driver not found",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    driver,
	})
}

func (h *httpService) findAllDriversHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

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

	query := model.DriverQueryInput{}
	if err := c.Bind(&query); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	drivers, total, err := h.driverRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.Errorf("Error querying drivers: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data: map[string]interface{}{
			"total":   total,
			"drivers": drivers,
		},
	})
}

func (h *httpService) createDriverHandler(c echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(c))

	var driver model.Driver

	if err := c.Bind(&driver); err != nil {
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

	if err := h.driverRepo.Create(c.Request().Context(), driver); err != nil {
		logger.Errorf("Error creating driver: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, response{
		Success: true,
		Data:    driver,
	})
}

func (h *httpService) updateDriverHandler(c echo.Context) error {
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

	var driver model.Driver

	if err := c.Bind(&driver); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	driver.ID = id

	if err := h.driverRepo.Update(c.Request().Context(), id, driver); err != nil {
		logger.Errorf("Error updating driver: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    driver,
	})
}

func (h *httpService) deleteDriverHandler(c echo.Context) error {
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

	if err := h.driverRepo.Delete(c.Request().Context(), id); err != nil {
		logger.Errorf("Error deleting driver: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
	})
}
