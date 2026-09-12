package port

import "context"

type SignInGoogleInput struct {
	Code string `json:"code" validate:"required"`
}

type SignInGoogleUseCase interface {
	Execute(ctx context.Context, input SignInGoogleInput) (*SignInOutput, error)
}
