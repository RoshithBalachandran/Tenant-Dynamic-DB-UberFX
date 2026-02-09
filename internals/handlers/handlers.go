package handlers

import (
	"tenant-Dynamin-DB/internals/config"
	"tenant-Dynamin-DB/internals/database"
	"tenant-Dynamin-DB/internals/models"
	"tenant-Dynamin-DB/internals/service"
	"tenant-Dynamin-DB/internals/token"

	"github.com/gofiber/fiber/v2"
)

type Handlers struct {
	cfg *config.Config
	srv *service.UserService
}

func NewHandlers(cfg *config.Config, srv *service.UserService) *Handlers {
	return &Handlers{cfg: cfg, srv: srv}
}
func (h *Handlers) Registration(c *fiber.Ctx) error {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Tenant   string `json:"tenant"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	db, err := database.ConnectDB(h.cfg, body.Tenant)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "migration failed"})
	}

	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: body.Password, // hash in service
		Tenant:   body.Tenant,
	}

	if err := h.srv.Register(db, &user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": "registered",
		"user":    user.Name,
	})
}

func (h *Handlers) LoginRequest(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Tenant   string `json:"tenant"`
	}
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": "invalid structure"})
	}
	db, err := database.ConnectDB(h.cfg, req.Tenant)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := h.srv.Login(db, req.Email, req.Password, req.Tenant)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	tk, err := token.GenerateToken(h.cfg.JWT_SECRET, user.ID, user.Tenant)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "token generation failed",
		})
	}

	return c.JSON(fiber.Map{
		"success": "login successful",
		"user":    user.Name,
		"token":   tk,
	})
}

func (h *Handlers) ListAll(c *fiber.Ctx) error {
	tenant, ok := c.Locals("tenant").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"Error": "unauthorzed"})
	}
	db, err := database.ConnectDB(h.cfg, tenant)
	if err != nil {
		return c.Status(fiber.StatusExpectationFailed).JSON(fiber.Map{"error": err.Error()})
	}
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"sucess": "user data fetched sucessfull", "user": users})
}

func (h *Handlers) Profile(c *fiber.Ctx) error {
	id := c.Locals("user_id")
	tena := c.Locals("tenant")

	if id == nil || tena == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "unauthorized"})
	}
	userid, ok := id.(uint)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user_id"})
	}
	tenant, ok := tena.(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tenant code"})
	}
	db, err := database.ConnectDB(h.cfg, tenant)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := h.srv.GetProfile(db, userid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{
		"id":     user.ID,
		"name":   user.Name,
		"email":  user.Email,
		"tenant": user.Tenant,
	})
}
