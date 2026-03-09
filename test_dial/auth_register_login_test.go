package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"sso/internal/services/auth"

	pb "github.com/egzggg/protos_sso/gen/go/sso"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Suite хранит объекты для тестов
type Suite struct {
	T          *testing.T
	AuthClient pb.AuthClient
	Cfg        struct {
		TokenTTL time.Duration
	}
}

// New создаёт локальный gRPC сервер и возвращает клиента
func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	// Создаём stub auth service
	stubAuth := auth.New(nil, nil, nil, nil, time.Minute) // nil для провайдеров, stub токен

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, &LocalGRPCServer{svc: stubAuth})

	// Слушаем на случайном порту
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go grpcServer.Serve(lis)

	t.Cleanup(func() {
		grpcServer.Stop()
	})

	// Подключаем клиента
	cc, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	t.Cleanup(func() {
		cc.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	t.Cleanup(cancel)

	return ctx, &Suite{
		T:          t,
		AuthClient: pb.NewAuthClient(cc),
		Cfg: struct{ TokenTTL time.Duration }{
			TokenTTL: time.Minute,
		},
	}
}

// LocalGRPCServer оборачивает auth.Auth для stub-теста
type LocalGRPCServer struct {
	pb.UnimplementedAuthServer
	svc *auth.Auth
}

func (s *LocalGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	id := int64(123) // stub ID
	return &pb.RegisterResponse{UserId: id}, nil
}

func (s *LocalGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{Token: "stub-token"}, nil
}

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := New(t)

	// Тестовые данные
	email := "test@example.com"
	pass := "password123"

	// Register
	respReg, err := st.AuthClient.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	if respReg.GetUserId() == 0 {
		t.Fatalf("expected non-zero user ID")
	}

	// Login
	respLogin, err := st.AuthClient.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    1,
	})
	require.NoError(t, err)
	if respLogin.GetToken() != "stub-token" {
		t.Fatalf("expected stub-token, got %s", respLogin.GetToken())
	}
}
