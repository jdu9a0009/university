package auth

import (
	"net/http"
	"project/foundation/web"
	"project/internal/commands"
	"project/internal/repository/postgres/user"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Controller represents the controller for authentication operations.
type Controller struct {
	user User
}

// NewController creates a new authentication controller.
func NewController(user User) *Controller {
	return &Controller{user: user}
}

// SignIn handles the sign-in operation.
func (uc Controller) SignIn(c *web.Context) error {
	var data user.SignInRequest

	err := c.BindFunc(&data, "Username", "Password")
	if err != nil {
		return c.RespondError(err)
	}

	detail, err := uc.user.GetByUsername(c.Ctx, data.Username)
	if err != nil {
		return c.RespondError(err)
	}

	if detail.Password == nil {
		return c.RespondError(&web.Error{
			Err:    errors.New("user not found"),
			Status: http.StatusNotFound,
		})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*detail.Password), []byte(data.Password)); err != nil {
		return c.RespondError(web.NewRequestError(errors.New("incorrect password"), http.StatusBadRequest))
	}

	accessToken, refreshToken, err := commands.GenToken(user.AuthClaims{
		ID:   detail.ID,
		Role: *detail.Role,
	}, "./private.pem")

	if err != nil {
		return c.RespondError(err)
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"data": map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
		"error": nil,
	}, http.StatusOK)
}

// Refresh handles the refresh token operation.
func (uc Controller) Refresh(c *web.Context) error {
	var data user.RefreshRequest

	err := c.BindFunc(&data, "AccessToken", "RefreshToken")
	if err != nil {
		return c.RespondError(err)
	}

	// Parse the incoming tokens
	_, refreshTokenClaims, err := commands.VerifyTokens(data.AccessToken, data.RefreshToken, "./private.pem")
	if err != nil {
		return c.RespondError(web.NewRequestError(err, http.StatusUnauthorized))
	}
	// Generate new tokens
	userClaims := user.AuthClaims{
		ID:   refreshTokenClaims.UserId,
		Role: refreshTokenClaims.Role,
	}

	accessToken, refreshToken, err := commands.GenToken(userClaims, "./private.pem")
	if err != nil {
		return c.RespondError(web.NewRequestError(errors.Wrap(err, "generating new tokens"), http.StatusInternalServerError))
	}

	return c.Respond(map[string]interface{}{
		"status": true,
		"data": map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
		"error": nil,
	}, http.StatusOK)

}
