package server

import (
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router                *gin.Engine
	Config                *config.Config
	AuthHandler           *handler.AuthHandler
	ShopHandler           *handler.ShopHandler
	AddressHandler        *handler.AddressHandler
	TransactionHandler    *handler.TransactionHandler
	UserBankHandler       *handler.UserBankHandler
	ProductHandler        *handler.ProductHandler
	CategoryHandler       *handler.CategoryHandler
	StorefrontHandler     *handler.StorefrontHandler
	BankHandler           *handler.BankHandler
	ProductVariantHandler *handler.ProductVariantHandler
	ProductReviewHandler  *handler.ProductReviewHandler
	ShippingHandler       *handler.ShippingHandler
	CourierHandler        *handler.CourierHandler
	ProfileHandler        *handler.ProfileHandler
	SocketServerHandler   *handler.SocketServerHandler
	TestimonialHandler    *handler.TestimonialHandler
}

func (s *Server) SetHandler() {
	// socket server
	s.Router.GET("/socket.io/", s.SocketServerHandler.SocketServer)
	s.Router.POST("/socket.io/", s.SocketServerHandler.SocketServer)
	s.Router.Handle("WS", "/socket.io/", s.SocketServerHandler.SocketServer)
	s.Router.Handle("WSS", "/socket.io/", s.SocketServerHandler.SocketServer)

	s.Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-api-key, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	s.Router.Use(s.AuthHandler.ApiKey)

	s.Router.POST("/register", s.AuthHandler.Register)
	s.Router.POST("/login", s.AuthHandler.Login)
	s.Router.GET("/activation", s.AuthHandler.ActivateEmail)
	s.Router.POST("/validate_token", s.AuthHandler.ValidateToken)
	s.Router.POST("/change_password", s.AuthHandler.Authorizer, s.AuthHandler.ChangePassword)

	profile := s.Router.Group("/profile", s.AuthHandler.Authorizer)
	{
		profile.GET("", s.ProfileHandler.GetProfile)
		profile.PUT("", s.ProfileHandler.UpdateProfile)
		profile.PUT("/image", s.ProfileHandler.UpdateProfileImage)
	}

	shop := s.Router.Group("/shop")
	{
		shop.GET("", s.AuthHandler.Authorizer, s.ShopHandler.AllShop)
		shop.GET("/:id", s.ShopHandler.GetShop)
		// shop.POST("/:id/courier", s.AuthHandler.Authorizer, s.ShopCourier.CreateShopCourier)
		shop.POST("", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.CreateShop)
		shop.PUT("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.UpdateShop)
		shop.DELETE("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ShopHandler.DeleteShop)

		shop.GET("/:id/address", s.ShopHandler.GetMainAddress)
	}

	address := s.Router.Group("/address", s.AuthHandler.Authorizer)
	{
		s.Router.GET("/addresses", s.AuthHandler.Authorizer, s.AddressHandler.AllAddresses)
		address.GET("/:id", s.AddressHandler.GetAddress)
		address.POST("", s.AuthHandler.CheckActivatedAccount, s.AddressHandler.CreateAddress)
		address.PUT("/:id", s.AuthHandler.CheckActivatedAccount, s.AddressHandler.UpdateAddress)
		address.DELETE("/:id", s.AuthHandler.CheckActivatedAccount, s.AddressHandler.DeleteAddress)
	}

	transaction := s.Router.Group("/transaction")
	{
		s.Router.GET("transactions/:shop_id", s.AuthHandler.Authorizer, s.TransactionHandler.AllTransaction)
		transaction.GET("/:id", s.AuthHandler.Authorizer, s.TransactionHandler.GetTransaction)
		transaction.PUT("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.TransactionHandler.UpdateTransaction)
		s.Router.GET("/track_transaction/:transaction_code", s.TransactionHandler.GetTransaction)
		transaction.POST("", s.TransactionHandler.CreateTransaction)
		transaction.POST("/webhook", s.TransactionHandler.WebhookTransaction)
	}

	products := s.Router.Group("/product")
	{
		s.Router.GET("/products", s.AuthHandler.Authorizer, s.ProductHandler.AllProduct)
		s.Router.GET("/products/:shop_id", s.AuthHandler.Authorizer, s.ProductHandler.AllProduct)
		products.GET("/:id", s.AuthHandler.Authorizer, s.ProductHandler.GetProduct)
		s.Router.GET("/product_by_slug/:shop/:product", s.ProductHandler.GetProduct)
		products.POST("", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductHandler.CreateProduct)
		products.POST("/:id/uploadimage", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductHandler.UploadProductImage)
		products.PUT("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductHandler.UpdateProduct)
		products.DELETE("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductHandler.DeleteProduct)

		products.GET("/:id/variant", s.ProductVariantHandler.AllProductVariant)
		products.GET("/:id/variant/:varian_id", s.ProductVariantHandler.GetProductVariant)
		products.POST("/:id/variant", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductVariantHandler.CreateProductVariant)
		products.PUT("/:id/variant/:varian_id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductVariantHandler.UpdateProductVariant)
		products.DELETE("/:id/variant/:varian_id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.ProductVariantHandler.DeleteProductVariant)

		products.POST("/:id/review", s.ProductReviewHandler.CreateProductReview)
	}

	storefront := s.Router.Group("/storefront", s.AuthHandler.Authorizer)
	{
		s.Router.GET("/storefronts", s.AuthHandler.Authorizer, s.StorefrontHandler.AllStorefront)
		storefront.GET("/:id", s.StorefrontHandler.GetStorefront)
		storefront.POST("", s.AuthHandler.CheckActivatedAccount, s.StorefrontHandler.CreateStorefront)
		storefront.PUT("/:id", s.AuthHandler.CheckActivatedAccount, s.StorefrontHandler.UpdateStorefront)
		storefront.DELETE("/:id", s.AuthHandler.CheckActivatedAccount, s.StorefrontHandler.DeleteStorefront)
		storefront.GET("/:id/products", s.StorefrontHandler.GetProductsByStorefront)
		storefront.POST("/:id/product", s.StorefrontHandler.SetProductStorefront)
	}

	category := s.Router.Group("/category", s.AuthHandler.Authorizer)
	{
		category.GET("", s.CategoryHandler.AllCategory)
		category.GET("/:id", s.CategoryHandler.GetCategory)
	}

	bank := s.Router.Group("/bank")
	{
		s.Router.GET("/banks", s.BankHandler.AllBank)
		bank.GET("/:id", s.BankHandler.GetBank)
	}

	userbank := s.Router.Group("/userbank")
	{
		userbank.Use(cors.Default())
		s.Router.GET("/userbanks", s.AuthHandler.Authorizer, s.UserBankHandler.AllUserBank)
		s.Router.GET("/all_userbank", s.UserBankHandler.AllUserBank)
		userbank.GET("/:id", s.AuthHandler.Authorizer, s.UserBankHandler.GetUserBank)
		s.Router.GET("/userbank_by_shop_id/:shop_id", s.UserBankHandler.GetUserBankByShopID)
		userbank.POST("", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.UserBankHandler.CreateUserBank)
		userbank.PUT("/:id", s.AuthHandler.Authorizer, s.AuthHandler.CheckActivatedAccount, s.UserBankHandler.UpdateUserBank)
		userbank.POST("/webhook", s.UserBankHandler.WebhookUserBank)
		userbank.DELETE("/:id", s.AuthHandler.Authorizer, s.UserBankHandler.DeleteUserBank)
	}

	shipping := s.Router.Group("/shipping")
	{
		shipping.POST("/cost", s.ShippingHandler.Cost)
		shipping.POST("/waybill", s.ShippingHandler.Waybill)
	}

	// couriers := s.Router.Group("/courier")
	{
		s.Router.GET("/couriers", s.CourierHandler.AllCourier)
	}

	testimonial := s.Router.Group("/testimonial")
	{
		s.Router.GET("/testimonials", s.TestimonialHandler.AllTestimonials)
		testimonial.POST("", s.AuthHandler.Authorizer, s.TestimonialHandler.CreateTestimonial)
		testimonial.GET(":id", s.AuthHandler.Authorizer, s.TestimonialHandler.GetTestimonial)
	}
}

func (s *Server) Run() {
	s.SetHandler()
	s.Router.Run(":" + s.Config.Port)
}

func NewServer(
	config *config.Config,
	authhandler *handler.AuthHandler,
	shophandler *handler.ShopHandler,
	addresshandler *handler.AddressHandler,
	transactionhandler *handler.TransactionHandler,
	userbankhandler *handler.UserBankHandler,
	producthandler *handler.ProductHandler,
	categoryhandler *handler.CategoryHandler,
	storefronthandler *handler.StorefrontHandler,
	bankhandler *handler.BankHandler,
	productvarianthandler *handler.ProductVariantHandler,
	productReviewHandler *handler.ProductReviewHandler,
	shippingHandler *handler.ShippingHandler,
	courierHandler *handler.CourierHandler,
	profileHandler *handler.ProfileHandler,
	socketServerHandler *handler.SocketServerHandler,
	testimonialHandler *handler.TestimonialHandler,
) *Server {
	return &Server{
		Router:                gin.Default(),
		Config:                config,
		AuthHandler:           authhandler,
		ShopHandler:           shophandler,
		AddressHandler:        addresshandler,
		TransactionHandler:    transactionhandler,
		UserBankHandler:       userbankhandler,
		ProductHandler:        producthandler,
		CategoryHandler:       categoryhandler,
		StorefrontHandler:     storefronthandler,
		BankHandler:           bankhandler,
		ProductVariantHandler: productvarianthandler,
		ProductReviewHandler:  productReviewHandler,
		ShippingHandler:       shippingHandler,
		CourierHandler:        courierHandler,
		ProfileHandler:        profileHandler,
		SocketServerHandler:   socketServerHandler,
		TestimonialHandler:    testimonialHandler,
	}
}
