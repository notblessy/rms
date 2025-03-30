package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func (h *httpService) loginWithGoogleHandler(c echo.Context) error {
	logger := logrus.WithField("ctx", utils.Dump(c.Request().Context()))

	var authRequest model.AuthRequest

	if err := c.Bind(&authRequest); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, &response{
			Success: false,
			Message: err.Error(),
		})
	}

	authRequest.RequestOrigin = c.Request().Header.Get("Origin")

	auth, err := h.userRepo.Authenticate(c.Request().Context(), authRequest.Code, authRequest.RequestOrigin)
	if err != nil {
		logger.Errorf("Error verifying token: %v", err)
		return c.JSON(http.StatusUnauthorized, &response{
			Success: false,
			Message: "unauthorized",
		})
	}

	token, err := signJwtToken(auth.ID, auth.Name, auth.Role)
	if err != nil {
		logger.Errorf("Error signing token: %v", err)
		return c.JSON(http.StatusInternalServerError, &response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, &response{
		Success: true,
		Data: map[string]any{
			"token": token,
			"type":  "Bearer",
		},
	})
}

func (h *httpService) profileHandler(c echo.Context) error {
	logger := logrus.WithField("ctx", utils.Dump(c.Request().Context()))

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, &response{
			Success: false,
			Message: "unauthorized",
		})
	}

	user, err := h.userRepo.FindByID(c.Request().Context(), session.ID)
	if err != nil {
		logger.Errorf("Error querying user: %v", err)
		return c.JSON(http.StatusNotFound, &response{
			Success: false,
			Message: "user not found",
		})
	}

	return c.JSON(http.StatusOK, &response{
		Success: true,
		Data:    user,
	})
}

func (h *httpService) patchUserHandler(c echo.Context) error {
	logger := logrus.WithField("ctx", utils.Dump(c.Request().Context()))

	session, err := authSession(c)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return c.JSON(http.StatusUnauthorized, &response{
			Success: false,
			Message: "unauthorized",
		})
	}

	var user model.User

	if err := c.Bind(&user); err != nil {
		logger.Errorf("Error parsing request: %v", err)
		return c.JSON(http.StatusBadRequest, &response{
			Success: false,
			Message: err.Error(),
		})
	}

	err = h.userRepo.PatchUser(c.Request().Context(), session.ID, user)
	if err != nil {
		logger.Errorf("Error updating user: %v", err)
		return c.JSON(http.StatusInternalServerError, &response{
			Success: false,
			Message: "internal server error",
		})
	}

	return c.JSON(http.StatusOK, &response{
		Success: true,
	})
}
