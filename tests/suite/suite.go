package suite

import (
	"context"
	"log/slog"
	"net"
	"os"
	"sso/internal/app"
	"sso/internal/config"
	"strconv"
	"testing"

	ssov1 "github.com/egzggg/protos_sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

const (
	grpcHost = "localhost"
)

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	//t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml") // струкутура с инфорамацией о конфигурации приложения

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout) // контекст для аварийного завершения

	t.Cleanup(func() { // дефер для теста
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.Dial( // токен для соединения с сервером
		grpcAddress(cfg), // адрес хоста
		grpc.WithTransportCredentials(insecure.NewCredentials()), // тип соединения без шифрования
		grpc.WithBlock(), // задержка
	)

	if err != nil {
		t.Fatalf("Соединение не установлено: %v", err)
	}

	t.Cleanup(func() { // еще один дефер для теста
		cc.Close()
	})

	return ctx, &Suite{
		T:          t,                       // наш тест
		Cfg:        cfg,                     // конфиг структура
		AuthClient: ssov1.NewAuthClient(cc), // интерфейс для вызова функций грпс сервера
	}
}

func NewForTest(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	//t.Parallel()

	cfg := config.MustLoadByPath("../config/local.yaml") // струкутура с инфорамацией о конфигурации приложения

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout) // контекст для аварийного завершения

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	t.Cleanup(func() {
		application.GRPCSrv.Stop()
	})

	t.Cleanup(func() { // дефер для теста
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.Dial( // токен для соединения с сервером
		grpcAddress(cfg), // адрес хоста
		grpc.WithTransportCredentials(insecure.NewCredentials()), // тип соединения без шифрования
		grpc.WithBlock(), // задержка
	)

	if err != nil {
		t.Fatalf("Соединение не установлено: %v", err)
	}

	t.Cleanup(func() { // еще один дефер для теста
		cc.Close()
	})

	return ctx, &Suite{
		T:          t,                       // наш тест
		Cfg:        cfg,                     // конфиг структура
		AuthClient: ssov1.NewAuthClient(cc), // интерфейс для вызова функций грпс сервера
	}
}
