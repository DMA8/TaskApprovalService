package grpc

import (
	"context"

	"gitlab.com/g6834/team31/auth/pkg/grpc_auth"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	client grpc_auth.AuthClient
	conn   *grpc.ClientConn
}

// var EnableGRPCTracingDialOption = grpc.WithUnaryInterceptor(grpc.UnaryClientInterceptor(clientInterceptor))

// func clientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
// 	// trace current request w/ child span
// 	parentSpan, ok := trace.FromContext(ctx)
// 	span, _ := parentSpan.NewChild(method)
// 	defer span.Finish()

// 	// new metadata, or copy of existing
// 	md, ok := metadata.FromContext(ctx)
// 	if !ok {
// 		md = metadata.New(nil)
// 	} else {
// 		md = md.Copy()
// 	}

// 	// append trace header to context metadata
// 	// header specification: https://cloud.google.com/trace/docs/faq
// 	md[headerKey] = append(
// 		md[headerKey], fmt.Sprintf("%s/%d;o=1", span.TraceID(), 0),
// 	)
// 	ctx = metadata.NewContext(ctx, md)

// 	return invoker(ctx, method, req, reply, cc, opts...)
// }

func New(ctx context.Context, host, port string) (*AuthClient, error) {
	connStr := host + port
	// opts := []grpc.DialOption{
	// 	grpc.WithUnaryInterceptor(
	// 		grpc_opentracing.UnaryClientInterceptor(
	// 			grpc_opentracing.WithTracer(opentracing.GlobalTracer()),
	// 		),
	// 	),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// }

	conn, err := grpc.DialContext(
		ctx,
		connStr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	client := grpc_auth.NewAuthClient(conn)
	return &AuthClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *AuthClient) Stop() error {
	return c.conn.Close()
}

func (c *AuthClient) Validate(ctx context.Context, in JWTTokens) (ValidateResponse, error) {
	ctx, span := otel.Tracer("team31_tasks").Start(ctx, "grpc_auth")
	defer span.End()
	md := metadata.Pairs(
		"key1", "value1",
		"key2", "value2",
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	response, err := c.client.Validate(ctx, &grpc_auth.Credential{
		AccessToken:  in.Access,
		RefreshToken: in.Refresh,
	},
	)
	if err != nil {
		return ValidateResponse{}, err
	}
	return ValidateResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		Login:        response.Login,
		Success:      response.Success,
		IsUpdate:     response.IsUpdate,
	}, nil
}
