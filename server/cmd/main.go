package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"uptime-monitor/server/internal/background"
	"uptime-monitor/server/internal/config"
	"uptime-monitor/server/internal/handler"
	"uptime-monitor/server/internal/middleware"
	"uptime-monitor/server/internal/repository"
	"uptime-monitor/server/internal/service"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	configPath := flag.String("c", "config.yaml", "path to config file")
	flag.Parse()

	// 1. Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// 2. Initialize database
	if err := repository.Init(&cfg.Database, cfg.Server.AdminUsername, cfg.Server.AdminPassword); err != nil {
		log.Fatalf("init repository: %v", err)
	}

	// 3. Create services
	authSvc := service.NewAuthService(&cfg.Server)
	monitorSvc := service.NewMonitorService()
	agentSvc := service.NewAgentService(&cfg.Agent)
	notifySvc := service.NewNotifyService()
	alertSvc := service.NewAlertService(notifySvc)
	metricsSvc := service.NewMetricsService()

	// 4. Create handlers
	authHandler := handler.NewAuthHandler(authSvc)
	monitorHandler := handler.NewMonitorHandler(monitorSvc)
	agentHandler := handler.NewAgentHandler(agentSvc, alertSvc, metricsSvc)
	alertHandler := handler.NewAlertHandler(notifySvc)
	dashboardHandler := handler.NewDashboardHandler(metricsSvc)
	metricsHandler := handler.NewMetricsHandler(metricsSvc)

	// 5. Setup Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS middleware
	r.Use(corsMiddleware())

	// Public routes
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/v1/agent/register", agentHandler.Register)

	// Agent routes (Agent Token auth)
	agentGroup := r.Group("/api/v1/agent")
	agentGroup.Use(middleware.AgentAuthMiddleware())
	{
		agentGroup.GET("/tasks", agentHandler.GetTasks)
		agentGroup.POST("/report", agentHandler.Report)
		agentGroup.POST("/heartbeat", agentHandler.Heartbeat)
	}

	// Web API routes (JWT auth)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authSvc))
	{
		// Auth
		api.POST("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", authHandler.Me)

		// Monitors
		api.GET("/monitors", monitorHandler.List)
		api.POST("/monitors", monitorHandler.Create)
		api.GET("/monitors/:id", monitorHandler.Get)
		api.PUT("/monitors/:id", monitorHandler.Update)
		api.DELETE("/monitors/:id", monitorHandler.Delete)
		api.PATCH("/monitors/:id/enabled", monitorHandler.ToggleEnabled)
		api.GET("/monitors/:id/metrics", metricsHandler.GetMonitorMetrics)

		// Alerts
		api.GET("/alerts/channels", alertHandler.ListChannels)
		api.POST("/alerts/channels", alertHandler.CreateChannel)
		api.PUT("/alerts/channels/:id", alertHandler.UpdateChannel)
		api.DELETE("/alerts/channels/:id", alertHandler.DeleteChannel)
		api.POST("/alerts/channels/:id/test", alertHandler.TestChannel)
		api.GET("/alerts/history", alertHandler.ListHistory)

		// Dashboard
		api.GET("/dashboard/overview", dashboardHandler.Overview)

		// Agents
		api.GET("/agents", agentHandler.List)
		api.DELETE("/agents/:id", agentHandler.Delete)
	}

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metricsSvc.GetRegistry(), promhttp.HandlerOpts{})))

	// SPA static files — serve embedded web/dist for any unmatched route
	distFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("embed web dist: %v", err)
	}
	spaFS := http.FS(distFS)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Try to serve the file directly
		if path != "/" && !strings.HasPrefix(path, "/api/") {
			c.FileFromFS(path, spaFS)
			return
		}
		// Fall back to index.html for SPA routing
		c.FileFromFS("/index.html", spaFS)
	})

	// 6. Start background jobs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	background.StartJobs(ctx, cfg)

	// 7. Start HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// 8. Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	cancel() // stop background jobs

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}

// corsMiddleware adds CORS headers for development.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
