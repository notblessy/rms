package router

import (
	"github.com/labstack/echo/v4"
	"github.com/notblessy/rms/model"
	"gorm.io/gorm"
)

type httpService struct {
	db            *gorm.DB
	userRepo      model.UserRepository
	camperRepo    model.CamperRepository
	equipmentRepo model.EquipmentRepository
	driverRepo    model.DriverRepository
}

func NewHTTPService() *httpService {
	return &httpService{}
}

func (h *httpService) RegisterDB(db *gorm.DB) {
	h.db = db
}

func (h *httpService) RegisterUserRepository(u model.UserRepository) {
	h.userRepo = u
}

func (h *httpService) RegisterCamperRepository(c model.CamperRepository) {
	h.camperRepo = c
}

func (h *httpService) RegisterEquipmentRepository(e model.EquipmentRepository) {
	h.equipmentRepo = e
}

func (h *httpService) RegisterDriverRepository(d model.DriverRepository) {
	h.driverRepo = d
}

func (h *httpService) Routes(e *echo.Echo) {
	e.GET("/ping", h.ping)
	e.GET("/health", h.health)

	v1 := e.Group("/v1")
	v1.GET("/auth/google", h.loginWithGoogleHandler)

	publicCampers := v1.Group("/campers")
	publicCampers.GET("", h.findAllCampersHandler)
	publicCampers.GET("/:id", h.findCamperByIDHandler)

	v1.Use(NewJWTMiddleware().ValidateJWT)

	users := v1.Group("/users")
	users.GET("", h.findAllUserHandler)
	users.GET("/me", h.profileHandler)
	users.PATCH("", h.patchUserHandler)

	campers := v1.Group("/campers")
	campers.POST("", h.createCamperHandler)
	campers.PUT("/:id", h.updateCamperHandler)
	campers.DELETE("/:id", h.deleteCamperHandler)

	equipments := v1.Group("/equipments")
	equipments.GET("", h.findAllEquipmentHandler)
	equipments.GET("/:id", h.findEquipmentByIDHandler)
	equipments.POST("", h.createEquipmentHandler)
	equipments.PUT("/:id", h.updateEquipmentHandler)
	equipments.DELETE("/:id", h.deleteEquipmentHandler)

	drivers := v1.Group("/drivers")
	drivers.GET("", h.findAllDriversHandler)
	drivers.GET("/:id", h.findDriverByIDHandler)
	drivers.POST("", h.createDriverHandler)
	drivers.PUT("/:id", h.updateDriverHandler)
	drivers.DELETE("/:id", h.deleteDriverHandler)
}

func (h *httpService) ping(c echo.Context) error {
	return c.JSON(200, response{Data: "pong"})
}

func (h *httpService) health(c echo.Context) error {
	err := h.db.Raw("SELECT 1").Error
	if err != nil {
		return c.JSON(500, response{Message: err.Error()})
	}

	return c.JSON(200, "OK")
}
