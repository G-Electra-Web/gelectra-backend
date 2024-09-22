package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/gelectra/gelectra-backend/api"
	"github.com/gelectra/gelectra-backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

func main() {
	// create koanf
	k := koanf.New(".")

	cfgPath := flag.String("config", "config.toml", "config file for the app")
	flag.Parse()

	// Load config file
	if err := k.Load(file.Provider(*cfgPath), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	log := logger.New()

	//Initilize logger
	lgrOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if k.Bool("app.debug") {
		lgrOptions.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, lgrOptions))
	logger.Info("Loaded Config file")

	// Initialise DB
	database.Connect(k, logger)

	//Init Fiber
	app := fiber.New()

	// Initialize default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize default config (Assign the middleware to /metrics)
	app.Get("/metrics", monitor.New())
	// Initialize default config
	app.Use(log)
	api.InitRoutes(app)
	app.Listen(":3000")

}
