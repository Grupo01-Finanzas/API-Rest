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

	// Get host from environment or use default
	host := os.Getenv("HOST")
	if host == "" {
		host = cfg.ServerHost
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.ServerPort
	}

	// Initialize database connection
	db := cfg.DB

	// Migrate the database
	if err := migrateDB(db); err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	establishmentRepo := repository.NewEstablishmentRepository(db)
	clientRepo := repository.NewClientRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	productRepo := repository.NewProductRepository(db)
	creditAccountRepo := repository.NewCreditAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	lateFeeRepo := repository.NewLateFeeRepository(db)
	lateFeeRuleRepo := repository.NewLateFeeRuleRepository(db)
	installmentRepo := repository.NewInstallmentRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, establishmentRepo, cfg.JwtSecret)
	clientService := service.NewClientService(clientRepo, userRepo)
	establishmentService := service.NewEstablishmentService(establishmentRepo)
	adminService := service.NewAdminService(adminRepo, establishmentRepo, userRepo)
	productService := service.NewProductService(productRepo)
	creditAccountService := service.NewCreditAccountService(creditAccountRepo, transactionRepo, clientRepo, establishmentRepo, installmentRepo)
	transactionService := service.NewTransactionService(transactionRepo, creditAccountRepo)
	lateFeeService := service.NewLateFeeService(lateFeeRepo)
	lateFeeRuleService := service.NewLateFeeRuleService(lateFeeRuleRepo)
	installmentService := service.NewInstallmentService(installmentRepo)
	purchaseService := service.NewPurchaseService(userRepo, establishmentRepo, productRepo, creditAccountRepo, transactionRepo, installmentRepo, db)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	clientController := controller.NewClientController(clientService)
	establishmentController := controller.NewEstablishmentController(establishmentService)
	adminController := controller.NewAdminController(adminService)
	productController := controller.NewProductController(productService)
	creditAccountController := controller.NewCreditAccountController(creditAccountService)
	transactionController := controller.NewTransactionController(transactionService)
	lateFeeController := controller.NewLateFeeController(lateFeeService)
	lateFeeRuleController := controller.NewLateFeeRuleController(lateFeeRuleService)
	installmentController := controller.NewInstallmentController(installmentService)
	purchaseController := controller.NewPurchaseController(purchaseService)

	// Initialize Gin router
	router := gin.Default()

	// Change to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Recovery middleware
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(middleware.CorsMiddleware())

	// Swagger documentation
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Public routes
	publicRoutes := router.Group("/api/v1")
	{
		publicRoutes.POST("/register", authController.Register)
		publicRoutes.POST("/login", authController.Login)
		publicRoutes.POST("/refresh", authController.RefreshToken)
	}

	// Protected routes (require authentication)
	protectedRoutes := router.Group("/api/v1", middleware.AuthMiddleware(cfg.JwtSecret))
	{
		// Establishment routes
		protectedRoutes.POST("/establishments", adminController.RegisterEstablishment)
		protectedRoutes.PUT("/establishments/:id", establishmentController.UpdateEstablishment)
		protectedRoutes.DELETE("/establishments/:id", establishmentController.DeleteEstablishment)
		// Use a different parameter name for the client route
		protectedRoutes.PUT("/establishments/:id/clients/:client_id", establishmentController.AddClientToEstablishment)

		// Client routes
		protectedRoutes.GET("/clients", clientController.GetAllClients)
		protectedRoutes.GET("/clients/:id", clientController.GetClientByID)
		protectedRoutes.PUT("/clients/:id", clientController.UpdateClient)
		protectedRoutes.DELETE("/clients/:id", clientController.DeleteClient)

		// Admin routes
		protectedRoutes.GET("/admins", adminController.GetAllAdmins)
		protectedRoutes.GET("/admins/:id", adminController.GetAdminByID)
		protectedRoutes.PUT("/admins/:id", adminController.UpdateAdmin)
		protectedRoutes.DELETE("/admins/:id", adminController.DeleteAdmin)
		protectedRoutes.POST("/register-establishments", adminController.RegisterEstablishment)

		// Authentication route (reset password)
		protectedRoutes.POST("/reset-password", authController.ResetPassword)

		// Product routes
		protectedRoutes.POST("/products", productController.CreateProduct)
		protectedRoutes.GET("/products", productController.GetAllProducts)
		protectedRoutes.GET("/products/:id", productController.GetProductByID)
		protectedRoutes.GET("/establishments/:establishment_id/products", productController.GetProductsByEstablishmentID)
		protectedRoutes.PUT("/products/:id", productController.UpdateProduct)
		protectedRoutes.DELETE("/products/:id", productController.DeleteProduct)

		// Credit Account Routes
		protectedRoutes.POST("/credit-accounts", creditAccountController.CreateCreditAccount)
		protectedRoutes.GET("/credit-accounts/:id", creditAccountController.GetCreditAccountByID)
		protectedRoutes.PUT("/credit-accounts/:id", creditAccountController.UpdateCreditAccount)
		protectedRoutes.DELETE("/credit-accounts/:id", creditAccountController.DeleteCreditAccount)
		protectedRoutes.GET("/establishments/:establishment_id/credit-accounts", creditAccountController.GetCreditAccountsByEstablishmentID)
		protectedRoutes.GET("/clients/:id/credit-accounts", creditAccountController.GetCreditAccountsByClientID)
		protectedRoutes.POST("/establishments/:establishment_id/credit-accounts/apply-interest", creditAccountController.ApplyInterestToAllAccounts)
		protectedRoutes.POST("/establishments/:establishment_id/credit-accounts/apply-late-fees", creditAccountController.ApplyLateFeesToAllAccounts)
		protectedRoutes.GET("/establishments/:establishment_id/credit-accounts/debt-summary", creditAccountController.GetAdminDebtSummary)
		protectedRoutes.POST("/credit-accounts/:id/purchases", creditAccountController.ProcessPurchase)
		protectedRoutes.POST("/credit-accounts/:id/payments", creditAccountController.ProcessPayment)
		protectedRoutes.PUT("/credit-accounts/:id/clients/:id", creditAccountController.AssignCreditAccountToClient)

		// Credit Request Routes
		protectedRoutes.POST("/credit-requests", creditAccountController.CreateCreditRequest)
		protectedRoutes.GET("/credit-requests/:id", creditAccountController.GetCreditRequestByID)
		protectedRoutes.PUT("/credit-requests/:id/approve", creditAccountController.ApproveCreditRequest)
		protectedRoutes.PUT("/credit-requests/:id/reject", creditAccountController.RejectCreditRequest)
		protectedRoutes.GET("/establishments/:establishment_id/credit-requests/pending", creditAccountController.GetPendingCreditRequests)

		// Transaction Routes
		protectedRoutes.POST("/transactions", transactionController.CreateTransaction)
		protectedRoutes.GET("/transactions/:id", transactionController.GetTransactionByID)
		protectedRoutes.PUT("/transactions/:id", transactionController.UpdateTransaction)
		protectedRoutes.DELETE("/transactions/:id", transactionController.DeleteTransaction)
		protectedRoutes.GET("/credit-accounts/:id/transactions", transactionController.GetTransactionsByCreditAccountID)

		// Purchase Routes
		protectedRoutes.POST("/purchases", purchaseController.CreatePurchase)
		protectedRoutes.GET("/clients/:id/balance", purchaseController.GetClientBalance)
		protectedRoutes.GET("/clients/:id/transactions", purchaseController.GetClientTransactions)
		protectedRoutes.GET("/clients/:id/overdue-balance", purchaseController.GetClientOverdueBalance)
		//protectedRoutes.GET("/clients/:id/credit-accounts", purchaseController.GetClientCreditAccount)

		// Late Fee Routes
		protectedRoutes.POST("/late-fees", lateFeeController.CreateLateFee)
		protectedRoutes.GET("/late-fees/:id", lateFeeController.GetLateFeeByID)
		protectedRoutes.PUT("/late-fees/:id", lateFeeController.UpdateLateFee)
		protectedRoutes.DELETE("/late-fees/:id", lateFeeController.DeleteLateFee)
		protectedRoutes.GET("/credit-accounts/:id/late-fees", lateFeeController.GetLateFeesByCreditAccountID)

		// Late Fee Rule Routes
		protectedRoutes.POST("/late-fee-rules", lateFeeRuleController.CreateLateFeeRule)
		protectedRoutes.GET("/late-fee-rules/:id", lateFeeRuleController.GetLateFeeRuleByID)
		protectedRoutes.PUT("/late-fee-rules/:id", lateFeeRuleController.UpdateLateFeeRule)
		protectedRoutes.DELETE("/late-fee-rules/:id", lateFeeRuleController.DeleteLateFeeRule)
		protectedRoutes.GET("/late-fee-rules", lateFeeRuleController.GetAllLateFeeRules)
		protectedRoutes.GET("/establishments/:establishment_id/late-fee-rules", lateFeeRuleController.GetLateFeeRulesByEstablishmentID)

		// Installment Routes
		protectedRoutes.POST("/installments", installmentController.CreateInstallment)
		protectedRoutes.GET("/installments/:id", installmentController.GetInstallmentByID)
		protectedRoutes.PUT("/installments/:id", installmentController.UpdateInstallment)
		protectedRoutes.DELETE("/installments/:id", installmentController.DeleteInstallment)
		protectedRoutes.GET("/credit-accounts/:id/installments", installmentController.GetInstallmentsByCreditAccountID)
		protectedRoutes.GET("/credit-accounts/:id/installments/overdue", installmentController.GetOverdueInstallments)

	}

	// Start the server
	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

// Migrate the database tables
func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&entities.User{},
		&entities.Admin{},
		&entities.Client{},
		&entities.Establishment{},
		&entities.Product{},
		&entities.CreditAccount{},
		&entities.Transaction{},
		&entities.LateFee{},
		&entities.LateFeeRule{},
		&entities.Installment{},
		&entities.CreditAccountHistory{},
	)
}
