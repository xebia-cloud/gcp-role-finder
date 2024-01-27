/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"
	"xebia-gcloud/gcp-role-finder/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serves the Google Cloud Platform IAM roles query service",
	Long: `
Provides an quick an easy user interface to find the Google Cloud Platform
IAM role with the least number of permissions.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")

		handler, err := handlers.NewRoleHandler(ctx, roleRepository)
		if err != nil {
			return err
		}

		if err := handler.RefreshRoles(ctx); err != nil {
			return err
		}

		app := fiber.New()
		app.Use(limiter.New(limiter.Config{
			Max:               20,
			Expiration:        3 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}))
		app.Use(cors.New(cors.Config{
			AllowOrigins:  "*",
			ExposeHeaders: "Content-Range",
		}))
		app.Get("/roles", handler.List)
		app.Get("/roles/:id", handler.GetRoleByID)
		app.Get("/refresh", handler.Refresh)
		app.Static("/", "./website/dist")
		return app.Listen(fmt.Sprintf("%s:%d", host, port))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 8080, "to listen on")
	serveCmd.Flags().String("host", "", "to listen on")
}
