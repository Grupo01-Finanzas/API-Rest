package main

import (
	"ApiRestFinance/internal/config"
	"ApiRestFinance/internal/controller"
	"ApiRestFinance/internal/middleware"
	"ApiRestFinance/internal/model/entities"
	"ApiRestFinance/internal/repository"
	"ApiRestFinance/internal/service"

	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"

	_ "ApiRestFinance/docs" // Import swagger docs for documentation

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title Final Assignment Finance API Rest
// @version 1.0
// @description API for managing finances in small businesses.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading configuration: ", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	db := cfg.DB

	// Migrate the database
	if err := migrateDB(db); err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewClientRepository(db)
	establishmentRepo := repository.NewEstablishmentRepository(db)
	productRepo := repository.NewProductRepository(db)
	creditAccountRepo := repository.NewCreditAccountRepository(db, userRepo)
	transactionRepo := repository.NewTransactionRepository(db)
	installmentRepo := repository.NewInstallmentRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, establishmentRepo, cfg.JwtSecret)
	userService := service.NewUserService(userRepo, creditAccountRepo) // Initialize userService
	adminService := service.NewAdminService(establishmentRepo, userRepo)
	establishmentService := service.NewEstablishmentService(establishmentRepo, userRepo)
	productService := service.NewProductService(productRepo, establishmentRepo)
	creditAccountService := service.NewCreditAccountService(creditAccountRepo, transactionRepo, installmentRepo, clientRepo, establishmentRepo) // Update to use userRepo
	transactionService := service.NewTransactionService(transactionRepo, creditAccountRepo)
	installmentService := service.NewInstallmentService(installmentRepo)
	purchaseService := service.NewPurchaseService(userRepo, establishmentRepo, productRepo, creditAccountRepo, transactionRepo, installmentRepo)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService, adminService, creditAccountService) // Use the new UserController
	establishmentController := controller.NewEstablishmentController(establishmentService)
	productController := controller.NewProductController(productService)
	creditAccountController := controller.NewCreditAccountController(creditAccountService)
	transactionController := controller.NewTransactionController(transactionService)
	installmentController := controller.NewInstallmentController(installmentService)
	purchaseController := controller.NewPurchaseController(purchaseService)

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery())
	router.Use(middleware.CorsMiddleware())

	// Swagger documentation
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Public routes
	publicRoutes := router.Group("/api/v1")
	{
		publicRoutes.POST("/register", authController.RegisterAdmin)
		publicRoutes.POST("/login", authController.Login)
		publicRoutes.POST("/refresh", authController.RefreshToken)
	}

	// Protected routes (require authentication)
	protectedRoutes := router.Group("/api/v1", middleware.AuthMiddleware(cfg.JwtSecret))
	{
		// User routes
		protectedRoutes.POST("/clients", userController.CreateClient)
		protectedRoutes.GET("/users/:id", userController.GetUserByID)
		protectedRoutes.PUT("/users/:id", userController.UpdateUser)
		protectedRoutes.DELETE("/users/:id", userController.DeleteUser)
		protectedRoutes.GET("/admins/me", userController.GetAdminProfile)
		protectedRoutes.PUT("/admins/me", userController.UpdateAdminProfile)
		protectedRoutes.GET("/establishments/:establishmentID/clients", userController.GetClientsByEstablishmentID)
		protectedRoutes.POST("/users/:id/photo", userController.UploadUserPhoto)

		// Establishment routes
		protectedRoutes.GET("/establishments/me", establishmentController.GetEstablishment)
		protectedRoutes.PUT("/establishments/me", establishmentController.UpdateEstablishment)

		// Product routes
		protectedRoutes.POST("/products", productController.CreateProduct)
		protectedRoutes.GET("/products/:id", productController.GetProductByID)
		protectedRoutes.GET("/establishments/:establishmentID/products", productController.GetAllProductsByEstablishmentID)
		protectedRoutes.PUT("/products/:id", productController.UpdateProduct)
		protectedRoutes.DELETE("/products/:id", productController.DeleteProduct)

		// Credit Account Routes
		protectedRoutes.POST("/credit-accounts", creditAccountController.CreateCreditAccount)
		protectedRoutes.GET("/credit-accounts/:id", creditAccountController.GetCreditAccountByID)
		protectedRoutes.PUT("/clients/:clientID/credit-account", userController.UpdateClientCreditAccount)
		protectedRoutes.DELETE("/credit-accounts/:id", creditAccountController.DeleteCreditAccount)
		protectedRoutes.GET("/establishments/:establishmentID/credit-accounts", creditAccountController.GetCreditAccountsByEstablishmentID)
		protectedRoutes.GET("/clients/:clientID/credit-account", creditAccountController.GetCreditAccountByClientID)
		protectedRoutes.POST("/credit-accounts/:id/apply-interest", creditAccountController.ApplyInterestToAccount)
		protectedRoutes.POST("/credit-accounts/:id/apply-late-fee", creditAccountController.ApplyLateFeeToAccount)
		protectedRoutes.GET("/credit-accounts/overdue", creditAccountController.GetOverdueCreditAccounts)
		protectedRoutes.POST("/credit-accounts/:id/purchases", creditAccountController.ProcessPurchase)
		protectedRoutes.POST("/credit-accounts/:id/payments", creditAccountController.ProcessPayment)
		protectedRoutes.GET("/credit-accounts/debt-summary", creditAccountController.GetAdminDebtSummary)

		// Transaction Routes
		protectedRoutes.POST("/transactions", transactionController.CreateTransaction)
		protectedRoutes.GET("/transactions/:id", transactionController.GetTransactionByID)
		protectedRoutes.PUT("/transactions/:id", transactionController.UpdateTransaction)
		protectedRoutes.DELETE("/transactions/:id", transactionController.DeleteTransaction)
		protectedRoutes.GET("/credit-accounts/:id/transactions", transactionController.GetTransactionsByCreditAccountID)
		protectedRoutes.POST("/transactions/:id/confirm", transactionController.ConfirmPayment)

		// Purchase Routes
		protectedRoutes.POST("/purchases", purchaseController.CreatePurchase)
		protectedRoutes.GET("/clients/me/balance", purchaseController.GetClientBalance)
		protectedRoutes.GET("/clients/me/transactions", purchaseController.GetClientTransactions)
		protectedRoutes.GET("/clients/me/overdue-balance", purchaseController.GetClientOverdueBalance)
		protectedRoutes.GET("/clients/me/installments", purchaseController.GetClientInstallments)
		protectedRoutes.GET("/clients/me/credit-account", purchaseController.GetClientCreditAccount)
		protectedRoutes.GET("/clients/me/account-summary", purchaseController.GetClientAccountSummary)     // New endpoint
		protectedRoutes.GET("/clients/me/account-statement", purchaseController.GetClientAccountStatement) // New endpoint
		protectedRoutes.GET("/clients/me/account-statement/pdf", purchaseController.GetClientAccountStatementPDF)

		// Installment Routes
		protectedRoutes.POST("/installments", installmentController.CreateInstallment)
		protectedRoutes.GET("/installments/:id", installmentController.GetInstallmentByID)
		protectedRoutes.PUT("/installments/:id", installmentController.UpdateInstallment)
		protectedRoutes.DELETE("/installments/:id", installmentController.DeleteInstallment)
		protectedRoutes.GET("/credit-accounts/:id/installments", installmentController.GetInstallmentsByCreditAccountID)
		protectedRoutes.GET("/credit-accounts/:id/installments/overdue", installmentController.GetOverdueInstallments)

		// Authentication route (reset password)
		protectedRoutes.POST("/reset-password", authController.ResetPassword)
	}

	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

// Migrate the database tables
func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Establishment{},
		&entities.Product{},
		&entities.CreditAccount{},
		&entities.Transaction{},
		&entities.Installment{},
	)
}
