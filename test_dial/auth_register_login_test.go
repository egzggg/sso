package test_dial

import (
	"context"
	"sso/internal/storage/postgres"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgresConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := "postgres://postgres_auth:1234@localhost:5432/sso?sslmode=disable"

	storage, err := postgres.New(ctx, dbURL)
	require.NoError(t, err, "failed to connect to Postgres")

	// Проверяем соединение простым запросом SELECT 1
	var one int
	err = storage.Pool().QueryRow(ctx, "SELECT 1").Scan(&one)
	require.NoError(t, err, "failed to ping Postgres")
	require.Equal(t, 1, one)
}
