package ports

import (
	"context"

	"gitlab.com/g6834/team31/tasks/internal/adapters/grpc"
)

// ClientAuth Интерфейс grpc Клиента
type ClientAuth interface {
	Validate(ctx context.Context, in grpc.JWTTokens) (grpc.ValidateResponse, error)
}
