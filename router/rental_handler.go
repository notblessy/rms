package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findRentalByIDHandler(e echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(e))

	id := e.Param("id")

	_, err := authSession(e)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return e.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	rental, err := h.rentalRepo.FindByID(e.Request().Context(), id)
	if err != nil {
		logger.Errorf("Error querying rental: %v", err)
		return e.JSON(http.StatusNotFound, response{
			Success: false,
			Message: "rental not found",
		})
	}

	return e.JSON(http.StatusOK, response{
		Success: true,
		Data:    rental,
	})
}

func (h *httpService) findAllRentalHandler(e echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(e))

	session, err := authSession(e)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return e.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.IsCustomer() {
		logger.Errorf("User is not authorized to access this resource")
		return e.JSON(http.StatusForbidden, response{
			Success: false,
			Message: "forbidden",
		})
	}

	var rentalQuery model.RentalQueryInput

	if err := e.Bind(&rentalQuery); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return e.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	rentals, total, err := h.rentalRepo.FindAll(e.Request().Context(), rentalQuery)
	if err != nil {
		logger.Errorf("Error getting rentals: %v", err)
		return e.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: err.Error(),
		})
	}

	return e.JSON(http.StatusOK, withPaging(rentals, total, rentalQuery.PageOrDefault(), rentalQuery.SizeOrDefault()))
}

func (h *httpService) createRentalHandler(e echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(e))

	var rental model.RentalInput
	if err := e.Bind(&rental); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return e.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	session, err := authSession(e)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return e.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.ID == "" {
		logger.Errorf("User is not authorized to access this resource")
		return e.JSON(http.StatusForbidden, response{
			Success: false,
			Message: "forbidden",
		})
	}

	rental.CustomerID = session.ID

	err = h.rentalRepo.Create(e.Request().Context(), rental)
	if err != nil {
		logger.Errorf("Error creating rental: %v", err)
		return e.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: err.Error(),
		})
	}

	return e.JSON(http.StatusCreated, response{
		Success: true,
		Data:    rental,
	})
}

func (h *httpService) updateRentalHandler(e echo.Context) error {
	logger := logrus.WithField("context", utils.Dump(e))

	id := e.Param("id")

	var rental model.RentalInput
	if err := e.Bind(&rental); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return e.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	_, err := authSession(e)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return e.JSON(http.StatusUnauthorized, response{
			Success: false,
			Message: "unauthorized",
		})
	}

	err = h.rentalRepo.Update(e.Request().Context(), id, rental)
	if err != nil {
		logger.Errorf("Error updating rental: %v", err)
		return e.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: err.Error(),
		})
	}

	return e.JSON(http.StatusOK, response{
		Success: true,
		Data:    rental,
	})
}
