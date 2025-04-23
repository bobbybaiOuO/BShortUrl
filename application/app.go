package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bobbybaiOuO/BShortUrl/config"
	"github.com/bobbybaiOuO/BShortUrl/database"
	"github.com/bobbybaiOuO/BShortUrl/internal/api"
	"github.com/bobbybaiOuO/BShortUrl/internal/cache"
	"github.com/bobbybaiOuO/BShortUrl/internal/service"
	"github.com/bobbybaiOuO/BShortUrl/pkg/shortcode"
	"github.com/bobbybaiOuO/BShortUrl/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Application .
type Application struct {
	e *echo.Echo
	db *sql.DB
	redisClient *cache.RedisCache
	urlServer *service.URLService
	URLHandler *api.URLHandler
	cfg *config.Config
	ShortCodeGenerator *shortcode.ShortCode
}

// Init .
func (a *Application) Init(filePath string) error {
	cfg, err := config.LoadConfig(filePath)
	if err != nil {
		return fmt.Errorf("Loading config error: %w", err)
	}
	a.cfg = cfg

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return err
	}
	a.db = db

	redisClient, err := cache.NewRedisCache(cfg.Redis)
	if err != nil {
		return err
	}
	a.redisClient = redisClient

	a.ShortCodeGenerator = shortcode.NewShortCode(cfg.ShortCode.Length)

	a.urlServer = service.NewURLService(db, a.ShortCodeGenerator, cfg.App.DefaultDuration, redisClient, cfg.App.BaseURL)

	a.URLHandler = api.NewURLHandler(a.urlServer)

	e := echo.New()
	e.Server.WriteTimeout = cfg.Server.WriteTimeout
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	
	e.POST("/api/url", a.URLHandler.CreateURL)
	e.GET("/:code", a.URLHandler.RedirectURL)
	e.Validator = validator.NewCustomValidator()
	a.e = e
	return nil
}

// Run .
func (a *Application) Run() {
	go a.startServer()
	go a.cleanup()
	a.shutdown()
}

func (a *Application) startServer() {
	if err := a.e.Start(a.cfg.Server.Addr); err != nil {
		log.Println(err)
	}
}

func (a *Application) cleanup() {
	ticker := time.NewTicker(a.cfg.App.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := a.urlServer.DeleteURL(context.Background()); err != nil {
			log.Panicln(err)
		}
	}
}

func (a *Application) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	defer func() {
		if err := a.db.Close(); err != nil {
			log.Panicln(err)
		}
	}()
	
	defer func ()  {
		if err := a.redisClient.Close(); err != nil {
			log.Panicln(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := a.e.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}