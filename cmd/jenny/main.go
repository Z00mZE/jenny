package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Z00mZE/jenny/internal/app/jenny"
	"github.com/Z00mZE/jenny/pb/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, ctxCloser := context.WithCancel(context.Background())
	defer ctxCloser()

	// Регистрируем сервисы
	app, appError := jenny.NewApp(ctx)
	if appError != nil {
		log.Fatal(appError)
	}

	// Создаем TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Настройки gRPC сервера
	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(30 * time.Second),
		// Можно добавить middleware, интерцепторы и т.д.
	)
	service.RegisterApplicationServer(grpcServer, app)

	// Включаем reflection для тестирования (dev only)
	reflection.Register(grpcServer)

	log.Printf("gRPC server starting on %s", lis.Addr().String())

	// Graceful shutdown
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("Shutdown signal received...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	go func() {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("Forced shutdown: timeout exceeded")
			grpcServer.Stop()
		}
	}()

	grpcServer.GracefulStop()
	log.Println("Server stopped gracefully")
}
