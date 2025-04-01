package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findEquipmentByIDHandler(c echo.Context) error {
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

	equipment, err := h.equipmentRepo.FindByID(c.Request().Context(), id)
	if err != nil {
		logger.Errorf("Error querying equipment: %v", err)
		return c.JSON(http.StatusNotFound, response{
			Success: false,
			Message: "equipment not found",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    equipment,
	})
}

func (h *httpService) findAllEquipmentHandler(c echo.Context) error {
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

	query := model.EquipmentQueryInput{}
	if err := c.Bind(&query); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	equipments, total, err := h.equipmentRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.Errorf("Error querying equipments: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    withPaging(equipments, total, query.PageOrDefault(), query.SizeOrDefault()),
	})
}

func (h *httpService) createEquipmentHandler(c echo.Context) error {
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

	var equipment model.Equipment
	if err := c.Bind(&equipment); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := h.equipmentRepo.Create(c.Request().Context(), equipment); err != nil {
		logger.Errorf("Error creating equipment: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, response{
		Success: true,
		Data:    equipment,
	})
}

func (h *httpService) updateEquipmentHandler(c echo.Context) error {
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

	var equipment model.Equipment
	if err := c.Bind(&equipment); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := h.equipmentRepo.Update(c.Request().Context(), id, equipment); err != nil {
		logger.Errorf("Error updating equipment: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, response{
		Success: true,
		Data:    equipment,
	})
}

func (h *httpService) deleteEquipmentHandler(c echo.Context) error {
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

	if err := h.equipmentRepo.Delete(c.Request().Context(), id); err != nil {
		logger.Errorf("Error deleting equipment: %v", err)
		return c.JSON(http.StatusInternalServerError, response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusNoContent, response{
		Success: true,
	})
}
