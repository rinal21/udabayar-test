package cmd

import (
	"udabayar-go-api-di/config"
	"udabayar-go-api-di/database"
	"udabayar-go-api-di/handler"
	"udabayar-go-api-di/helper"
	"udabayar-go-api-di/server"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	// server
	container.Provide(server.NewServer)

	// config
	container.Provide(config.NewConfig)
	container.Provide(config.NewDatabaseConfig)
	container.Provide(config.NewEmailConfig)
	container.Provide(config.NewTransactionConfig)

	// database
	container.Provide(database.ConnectDatabase)
	container.Provide(database.DatabaseMigration)

	// helper
	container.Provide(helper.NewJWTHelper)
	container.Provide(helper.NewBcryptHelper)
	container.Provide(helper.NewEmailHelper)
	container.Provide(helper.NewUploadHelper)
	container.Provide(helper.NewTimeHelper)
	container.Provide(helper.NewUnitHelper)

	// handler
	container.Provide(handler.NewAuthHandler)
	container.Provide(handler.NewProfileHandler)
	container.Provide(handler.NewShopHandler)
	container.Provide(handler.NewAddressHandler)
	container.Provide(handler.NewTransactionHandler)
	container.Provide(handler.NewUserBankHandler)
	container.Provide(handler.NewProductHandler)
	container.Provide(handler.NewCategoryHandler)
	container.Provide(handler.NewStorefrontHandler)
	container.Provide(handler.NewBankHandler)
	container.Provide(handler.NewUserBankHandler)
	container.Provide(handler.NewProductVariantHandler)
	container.Provide(handler.NewProductReviewHandler)
	container.Provide(handler.NewShippingHandler)
	container.Provide(handler.NewCourierHandler)
	container.Provide(handler.NewTestimonialHandler)

	// socket
	container.Provide(handler.NewSocketServerHandler)

	return container
}

func Serve() {
	container := BuildContainer()

	err := container.Invoke(func(
		server *server.Server,
	) {
		server.Run()
	})

	if err != nil {
		panic(err)
	}
}

func MigrateDatabase() {
	container := BuildContainer()

	err := container.Invoke(func(
		database *database.Database,
	) {
		database.Migrate()
	})

	if err != nil {
		panic(err)
	}
}

func SeedPrimaryData() {
	container := BuildContainer()

	err := container.Invoke(func(
		database *database.Database,
	) {
		database.Seeder()
	})

	if err != nil {
		panic(err)
	}
}

func SeedDummyData() {
	container := BuildContainer()

	err := container.Invoke(func(
		database *database.Database,
	) {
		database.Seeder()
	})

	if err != nil {
		panic(err)
	}
}

func MigrateEmailTemplate() {
	container := BuildContainer()

	err := container.Invoke(func(
		email *helper.Email,
	) {
		email.CreateTemplate()
	})

	if err != nil {
		panic(err)
	}
}
