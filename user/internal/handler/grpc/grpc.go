package auth

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"wuphf.com/user/gen"
	"wuphf.com/user/internal/controller/user"
	"wuphf.com/user/pkg/auth"
)

// Handler defines a user gRPC handler
type Handler struct {
	gen.UnimplementedAuthServiceServer
	ctrl *user.Controller
}

// New creates a new user gRPC handler
func New(ctrl *user.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// ValidateToken validates a JWT token
func (h *Handler) ValidateToken(ctx context.Context, req *gen.TokenRequest) (*gen.TokenResponse, error) {
	tokenString := req.GetToken()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return auth.JWTKey, nil
	})
	if err != nil {
		return &gen.TokenResponse{Valid: false}, nil
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["userID"].(string)

			user, err := h.ctrl.Get(ctx, userID)
			if err != nil {
				return &gen.TokenResponse{Valid: false}, nil
			}

			return &gen.TokenResponse{Valid: true, User: &gen.User{Id: user.ID, Email: user.Email}}, nil
		}
	}

	return &gen.TokenResponse{Valid: false}, nil
}
