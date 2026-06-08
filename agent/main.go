package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := flag.String("c", "agent.yaml", "配置文件路径")
	flag.Parse()

	// Load configuration.
	cfg, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// Register with server (or load cached token).
	token, err := Register(cfg)
	if err != nil {
		log.Fatalf("Agent 注册失败: %v", err)
	}

	// Shared result queue.
	results := NewResultQueue()

	// Context for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())

	// Start scheduler and reporter in background goroutines.
	scheduler := NewScheduler(cfg, token, results)
	reporter := NewReporter(cfg, token, results)

	go scheduler.Start(ctx)
	go reporter.Start(ctx)

	log.Printf("Agent [%s] 已启动，等待任务...", cfg.Agent.Name)

	// Wait for shutdown signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("收到信号 %v，开始关闭...", sig)

	// Cancel context to stop all goroutines.
	cancel()

	// Wait for final result upload.
	time.Sleep(10 * time.Second)

	log.Println("Agent 已关闭")
}
