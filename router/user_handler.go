package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) findAllUserHandler(c echo.Context) error {
	logger := logrus.WithField("ctx", utils.Dump(c.Request().Context()))

	var query model.UserQueryInput

	if err := c.Bind(&query); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, &response{
			Success: false,
			Message: err.Error(),
		})
	}

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, &response{
			Success: false,
			Message: "unauthorized",
		})
	}

	if session.IsCustomer() {
		logger.Errorf("User is not authorized to access this resource")
		return c.JSON(http.StatusForbidden, &response{
			Success: false,
			Message: "forbidden",
		})
	}

	users, total, err := h.userRepo.FindAll(c.Request().Context(), query)
	if err != nil {
		logger.Errorf("Error getting users: %v", err)
		return c.JSON(http.StatusInternalServerError, &response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, withPaging(users, total, query.PageOrDefault(), query.SizeOrDefault()))
}
