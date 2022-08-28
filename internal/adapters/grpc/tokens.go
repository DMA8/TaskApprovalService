package grpc

type JWTTokens struct {
	Access  string
	Refresh string
}

type ValidateResponse struct {
	Login        string
	AccessToken  string
	RefreshToken string
	Success      bool
	IsUpdate     bool
}
