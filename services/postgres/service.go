package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init(logger *zap.Logger) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), viper.GetString("POSTGRES_CONN"))
	if err != nil {
		logger.Fatal(fmt.Sprintf("error while connecting to PSQL: %v", conn))
	}
	const maxConnAttempts = 3

	var successConn bool
	for i := 0; i < maxConnAttempts; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err := conn.Ping(ctx)
		cancel()
		if err == nil {
			successConn = true
			break
		}
		logger.Warn(fmt.Sprintf("error while ping to PSQL: %v", err))
	}
	if !successConn {
		logger.Fatal("can't establish connection to PSQL")
	}
	logger.Info("PSQL connected successfully")
	return conn
}
