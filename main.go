package main

import (
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/notblessy/rms/db"
	"github.com/notblessy/rms/repository"
	"github.com/notblessy/rms/router"
	"github.com/notblessy/rms/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("cannot load .env file")
	}

	postgres := db.NewPostgres()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Path",
		},
	}))
	e.Use(middleware.CORS())
	e.Validator = &utils.Ghost{Validator: validator.New()}

	userRepo := repository.NewUserRepository(postgres)
	camperRepo := repository.NewCamperRepository(postgres)
	equipmentRepo := repository.NewEquipmentRepository(postgres)
	driverRepo := repository.NewDriverRepository(postgres)

	httpService := router.NewHTTPService()
	httpService.RegisterDB(postgres)
	httpService.RegisterUserRepository(userRepo)
	httpService.RegisterCamperRepository(camperRepo)
	httpService.RegisterEquipmentRepository(equipmentRepo)
	httpService.RegisterDriverRepository(driverRepo)

	httpService.Routes(e)

	e.Logger.Fatal(e.Start(":3500"))
}
