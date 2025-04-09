package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	errInvalidRequest string = "not implemented"
)

type AuthController struct {
	authHandler AuthHandler
}

func NewAuthController(authHandler AuthHandler) *AuthController {
	return &AuthController{authHandler: authHandler}
}

// AuthHandler defines the interface for auth-related business logic
type AuthHandler interface {
	DummyLogin(role string) (string, error)
	Register(email, password, role string) (*UserResponse, error)
	Login(email, password string) (string, error)
}

// UserResponse represents the user data returned by the handler
type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// DummyLogin godoc
// @Summary Получение тестового токена
// @Description Генерирует тестовый токен для указанной роли
// @Tags auth
// @Accept json
// @Produce json
// @Param input body DummyLoginRequest true "Роль пользователя"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse
// @Router /dummyLogin [post]
func (c *AuthController) DummyLogin(ctx *fiber.Ctx) error {
	var req DummyLoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: errInvalidRequest})
	}

	token, err := c.authHandler.DummyLogin(req.Role)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(TokenResponse{Token: token})
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные пользователя"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var req RegisterRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request"})
	}

	user, err := c.authHandler.Register(req.Email, req.Password, req.Role)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// Login godoc
// @Summary Авторизация пользователя
// @Description Авторизует пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Учетные данные"
// @Success 200 {object} TokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request"})
	}

	token, err := c.authHandler.Login(req.Email, req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(TokenResponse{Token: token})
}

// Request/Response structures
type DummyLoginRequest struct {
	Role string `json:"role" validate:"required,oneof=employee moderator"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=employee moderator"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// PVZController handles PVZ-related endpoints
type PVZController struct {
	pvzHandler PVZHandler
}

func NewPVZController(pvzHandler PVZHandler) *PVZController {
	return &PVZController{pvzHandler: pvzHandler}
}

type PVZHandler interface {
	CreatePVZ(city string, userID uuid.UUID) (*PVZResponse, error)
	GetPVZs(startDate, endDate string, page, limit int) ([]PVZWithReceptionsResponse, error)
	OpenReception(pvzID, userID uuid.UUID) (*ReceptionResponse, error)
	CloseLastReception(pvzID, userID uuid.UUID) (*ReceptionResponse, error)
	DeleteLastProduct(pvzID, userID uuid.UUID) error
}

// PVZResponse represents PVZ data returned by the handler
type PVZResponse struct {
	ID               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

// CreatePVZ godoc
// @Summary Создание ПВЗ (только для модераторов)
// @Description Создает новый пункт выдачи заказов
// @Tags pvz
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param input body PVZRequest true "Данные ПВЗ"
// @Success 201 {object} PVZResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /pvz [post]
func (c *PVZController) CreatePVZ(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	var req PVZRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request"})
	}

	pvz, err := c.pvzHandler.CreatePVZ(req.City, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(pvz)
}

// GetPVZs godoc
// @Summary Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией
// @Description Возвращает список ПВЗ с информацией о приемках и товарах
// @Tags pvz
// @Security bearerAuth
// @Produce json
// @Param startDate query string false "Начальная дата диапазона"
// @Param endDate query string false "Конечная дата диапазона"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Success 200 {array} PVZWithReceptionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /pvz [get]
func (c *PVZController) GetPVZs(ctx *fiber.Ctx) error {
	_, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	page := ctx.QueryInt("page", 1)
	limit := ctx.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 30 {
		limit = 10
	}

	pvzs, err := c.pvzHandler.GetPVZs(startDate, endDate, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(pvzs)
}

// OpenReception godoc
// @Summary Открытие новой приемки товаров в рамках ПВЗ
// @Description Открывает новую приемку товаров для указанного ПВЗ
// @Tags pvz
// @Security bearerAuth
// @Produce json
// @Param pvzId path string true "ID ПВЗ"
// @Success 200 {object} ReceptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /pvz/{pvzId}/open_reception [post]
func (c *PVZController) OpenReception(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	pvzID, err := uuid.Parse(ctx.Params("pvzId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid PVZ ID"})
	}

	reception, err := c.pvzHandler.OpenReception(pvzID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(reception)
}

// CloseLastReception godoc
// @Summary Закрытие последней открытой приемки товаров в рамках ПВЗ
// @Description Закрывает последнюю открытую приемку товаров для указанного ПВЗ
// @Tags pvz
// @Security bearerAuth
// @Produce json
// @Param pvzId path string true "ID ПВЗ"
// @Success 200 {object} ReceptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /pvz/{pvzId}/close_last_reception [post]
func (c *PVZController) CloseLastReception(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	pvzID, err := uuid.Parse(ctx.Params("pvzId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid PVZ ID"})
	}

	reception, err := c.pvzHandler.CloseLastReception(pvzID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.JSON(reception)
}

// DeleteLastProduct godoc
// @Summary Удаление последнего добавленного товара из текущей приемки
// @Description Удаляет последний добавленный товар из текущей приемки (LIFO)
// @Tags pvz
// @Security bearerAuth
// @Produce json
// @Param pvzId path string true "ID ПВЗ"
// @Success 200
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /pvz/{pvzId}/delete_last_product [post]
func (c *PVZController) DeleteLastProduct(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	pvzID, err := uuid.Parse(ctx.Params("pvzId"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid PVZ ID"})
	}

	err = c.pvzHandler.DeleteLastProduct(pvzID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// PVZRequest represents PVZ creation request
type PVZRequest struct {
	City string `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
}

// PVZWithReceptionsResponse represents PVZ with receptions and products
type PVZWithReceptionsResponse struct {
	PVZ        PVZResponse             `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception ReceptionResponse `json:"reception"`
	Products  []ProductResponse `json:"products"`
}

type ReceptionResponse struct {
	ID       string `json:"id"`
	DateTime string `json:"dateTime"`
	PVZID    string `json:"pvzId"`
	Status   string `json:"status"`
}

type ProductResponse struct {
	ID          string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionID string `json:"receptionId"`
}

// ProductController handles product-related endpoints
type ProductController struct {
	productHandler ProductHandler
}

func NewProductController(productHandler ProductHandler) *ProductController {
	return &ProductController{productHandler: productHandler}
}

type ProductHandler interface {
	CreateProduct(productType string, pvzID, userID uuid.UUID) (*ProductResponse, error)
}

// CreateProduct godoc
// @Summary Добавление товара в текущую приемку
// @Description Добавляет товар в текущую открытую приемку для указанного ПВЗ
// @Tags products
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param input body ProductRequest true "Данные товара"
// @Success 201 {object} ProductResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /products [post]
func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	var req ProductRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request"})
	}

	pvzID, err := uuid.Parse(req.PVZID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid PVZ ID"})
	}

	product, err := c.productHandler.CreateProduct(req.Type, pvzID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(product)
}

// ProductRequest represents product creation request
type ProductRequest struct {
	Type  string `json:"type" validate:"required,oneof=электроника одежда обувь"`
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}

// ReceptionController handles reception-related endpoints
type ReceptionController struct {
	receptionHandler ReceptionHandler
}

func NewReceptionController(receptionHandler ReceptionHandler) *ReceptionController {
	return &ReceptionController{receptionHandler: receptionHandler}
}

type ReceptionHandler interface {
	CreateReception(pvzID, userID uuid.UUID) (*ReceptionResponse, error)
}

// CreateReception godoc
// @Summary Создание новой приемки товаров
// @Description Создает новую приемку товаров для указанного ПВЗ
// @Tags receptions
// @Security bearerAuth
// @Accept json
// @Produce json
// @Param input body ReceptionRequest true "Данные приемки"
// @Success 201 {object} ReceptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /receptions [post]
func (c *ReceptionController) CreateReception(ctx *fiber.Ctx) error {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: "Access denied"})
	}

	var req ReceptionRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid request"})
	}

	pvzID, err := uuid.Parse(req.PVZID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: "Invalid PVZ ID"})
	}

	reception, err := c.receptionHandler.CreateReception(pvzID, userID)
	if err != nil {
		if err.Error() == "access denied" {
			return ctx.Status(fiber.StatusForbidden).JSON(ErrorResponse{Message: err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Message: err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(reception)
}

// ReceptionRequest represents reception creation request
type ReceptionRequest struct {
	PVZID string `json:"pvzId" validate:"required,uuid4"`
}

// Helper function to get user ID from context (set by auth middleware)
func getUserIDFromContext(ctx *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := ctx.Locals("userID").(string)
	return uuid.Parse(userIDStr)
}
