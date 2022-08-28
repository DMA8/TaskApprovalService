package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"gitlab.com/g6834/team31/tasks/internal/adapters/grpc"
	"gitlab.com/g6834/team31/tasks/internal/config"

	uuid "github.com/satori/go.uuid"
)

type ctxKey int
type userLogin string

const (
	ridKey ctxKey    = 0
	usLog  userLogin = "userLogin"
)

var (
	ErrBadCredential = errors.New("bad credentials")
)

func (s *Server) ValidateToken(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		newCtx, span := tracer.Start(r.Context(), "middleware s.validateToken")
		defer span.End()
		cfg := config.NewConfig().HTTP
		access, err := r.Cookie(cfg.AccessCookieName)
		if err != nil {
			s.logger.Debug().Err(err).Msgf("s.ValidateToken bad access token")
			WriteAnswer(rw, http.StatusForbidden, "bad access", s.logger)
			return
		}
		refresh, err := r.Cookie(cfg.RefreshCookieName)
		if err != nil {
			s.logger.Debug().Err(err).Msgf("s.ValidateToken bad refresh token")
			WriteAnswer(rw, http.StatusForbidden, "bad refresh", s.logger)
			return
		}
		ctx := newCtx
		credential, err := s.AuthClient.Validate(ctx, grpc.JWTTokens{
			Access:  access.Value,
			Refresh: refresh.Value,
		})
		if err != nil {
			s.logger.Debug().Err(err).Msgf("s.ValidateToken auth couldn't validate cookies")
			WriteAnswer(rw, http.StatusForbidden, ErrBadCredential.Error(), s.logger)
			return
		}
		ctx = context.WithValue(ctx, usLog, credential.Login)
		if credential.IsUpdate {
			rw.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s", cfg.AccessCookieName, credential.AccessToken))
			rw.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s", cfg.RefreshCookieName, credential.RefreshToken))
		}
		s.logger.Debug().Msgf("s.ValidateToken auth succeed")
		next.ServeHTTP(rw, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewV4().String()
		}
		ctx := context.WithValue(r.Context(), ridKey, rid)
		w.Header().Add("X-Request-ID", rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(ridKey).(string)
}
