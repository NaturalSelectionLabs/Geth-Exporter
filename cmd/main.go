package cmd

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/naturalselectionlabs/geth-exporter/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	configFile string
	debug      bool
)

var rootCmd = cobra.Command{
	Use:   "geth-exporter",
	Short: "geth-exporter",
	Long:  "This is an application exporter geth json rpc.",
	RunE: func(cmd *cobra.Command, args []string) error {

		zap.L().Info("Loading config file", zap.String("file", configFile))
		cfg, err := config.LoadConfig(configFile)
		if err != nil {
			zap.L().Error("Error loading config", zap.String("file", configFile), zap.Error(err))
			return err
		}
		configJSON, err := json.Marshal(cfg)
		if err != nil {
			zap.L().Error("Failed to marshal config to JSON", zap.String("file", configFile), zap.Error(err))
			return err
		}
		zap.L().Info("Loaded config file", zap.String("config", string(configJSON)))

		server := echo.New()

		server.HidePort = !debug
		server.HideBanner = !debug

		server.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
		server.GET("/probe", func(c echo.Context) error {
			return probeHandler(c, cfg)
		})
		server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

		return server.Start(":8000")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		zap.L().Fatal("exec error:", zap.Error(err))
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug level")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")

	if debug {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
}
