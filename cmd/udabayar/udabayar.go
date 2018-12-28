package cmd

import (
	"fmt"
	"udabayar-go-api-di/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var bigText = `
    ╭╮╱╭╮╭━━━╮╭━━━╮╭━━╮╱╭━━━╮╭╮╱╱╭╮╭━━━╮╭━━━╮╱╱╱╱╭━━━╮╭╮╱╱╱╭━━╮
    ┃┃╱┃┃╰╮╭╮┃┃╭━╮┃┃╭╮┃╱┃╭━╮┃┃╰╮╭╯┃┃╭━╮┃┃╭━╮┃╱╱╱╱┃╭━╮┃┃┃╱╱╱╰┫┣╯
    ┃┃╱┃┃╱┃┃┃┃┃┃╱┃┃┃╰╯╰╮┃┃╱┃┃╰╮╰╯╭╯┃┃╱┃┃┃╰━╯┃╱╱╱╱┃┃╱╰╯┃┃╱╱╱╱┃┃
    ┃┃╱┃┃╱┃┃┃┃┃╰━╯┃┃╭━╮┃┃╰━╯┃╱╰╮╭╯╱┃╰━╯┃┃╭╮╭╯╭━━╮┃┃╱╭╮┃┃╱╭╮╱┃┃
    ┃╰━╯┃╭╯╰╯┃┃╭━╮┃┃╰━╯┃┃╭━╮┃╱╱┃┃╱╱┃╭━╮┃┃┃┃╰╮╰━━╯┃╰━╯┃┃╰━╯┃╭┫┣╮
    ╰━━━╯╰━━━╯╰╯╱╰╯╰━━━╯╰╯╱╰╯╱╱╰╯╱╱╰╯╱╰╯╰╯╰━╯╱╱╱╱╰━━━╯╰━━━╯╰━━╯
`

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "udabayar",
		Short: "Udabayar CLI is a backend commands for Udabayar Rest API.",
		Long:  bigText + "\n Udabayar CLI is a backend commands for Udabayar Rest API.",
		Args:  cobra.MinimumNArgs(1),
	}

	cmdServe := &cobra.Command{
		Use:   "serve",
		Short: "Serve is a commands for serve Udabayar Rest API.",
		Long:  bigText + "\n Serve is a commands for serve Udabayar Rest API.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText+"\nStart Udabayar Rest API ...", "\n- Production Mode")

			config := config.NewConfig()
			cmd.Printf("\nStarted, serving HTTP on :%s", config.Port)

			gin.SetMode(gin.ReleaseMode)
			Serve()
		},
	}

	cmdServeProd := &cobra.Command{
		Use:   "production",
		Short: "Serve Udabayar Rest API on Production Mode.",
		Long:  bigText + "\n Serve Udabayar Rest API on Production Mode.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText+"\nStart Udabayar Rest API ...", "\n- Production Mode")

			config := config.NewConfig()
			cmd.Printf("\nStarted, serving HTTP on :%s", config.Port)

			gin.SetMode(gin.ReleaseMode)
			Serve()
		},
	}

	cmdServeDev := &cobra.Command{
		Use:   "development",
		Short: "Serve Udabayar Rest API on Development Mode.",
		Long:  bigText + "\n Serve Udabayar Rest API on Development Mode.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText+"\n Start Udabayar Rest API ...", "\n- Development Mode")
			Serve()
		},
	}

	cmdMigrate := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate is a command for migration before run this app.",
		Long:  bigText + "\n Migrate is a command for migration before run this app.",
	}

	cmdMigrateDB := &cobra.Command{
		Use:   "database",
		Short: "Migrate database of this app.",
		Long:  bigText + "\n Migrate database of this app.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText+"\n Migrate Database ...", "\n- Development Mode")
			MigrateDatabase()

			cmd.Println("Migrate Datbase successful :)")
		},
	}

	cmdMigrateEmail := &cobra.Command{
		Use:   "email",
		Short: "Migrate email template of this app to AWS SES (Simple Email Service).",
		Long:  bigText + "\n Migrate email template of this app to AWS SES (Simple Email Service).",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText + "\n Migrate Email Template to AWS SES (Simple Email Service) ...")
			MigrateEmailTemplate()

			cmd.Println("Migrate Email Template successful :)")
		},
	}

	cmdSeeder := &cobra.Command{
		Use:   "seeder",
		Short: "Seeds data to database.",
		Long:  bigText + "\n Seeds data to database.",
	}

	cmdSeederRun := &cobra.Command{
		Use:   "run",
		Short: "Seeds primary data to database.",
		Long:  bigText + "\n Seeds primary data to database.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText + "\n Seeding Primary Data to Database ...")
			SeedPrimaryData()

			cmd.Println("Seeder Primary Data successful :)")
		},
	}

	cmdSeederDummy := &cobra.Command{
		Use:   "dummy",
		Short: "Seeds dummy data to database.",
		Long:  bigText + "\n Seeds dummy data to database.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(bigText + "\n Seeding Dummy Data to Database ...")
			SeedDummyData()

			cmd.Println("Seeder Dummy Data successful :)")
		},
	}

	rootCmd.AddCommand(cmdServe, cmdMigrate, cmdSeeder)
	cmdServe.AddCommand(cmdServeProd, cmdServeDev)
	cmdMigrate.AddCommand(cmdMigrateDB, cmdMigrateEmail)
	cmdSeeder.AddCommand(cmdSeederRun, cmdSeederDummy)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
