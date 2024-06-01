package auth

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"project/foundation/web"
	"project/internal/commands"
	"project/internal/repository/postgres/user"
	"time"

	"github.com/golang-jwt/jwt"
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

	err := c.BindFunc(&data, "Login", "Password")
	if err != nil {
		return c.RespondError(err)
	}

	detail, err := uc.user.GetByLogin(c.Ctx, data.Username)
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
	accessClaims, err := parseToken(data.AccessToken, "./public.pem")
	if err != nil {
		return c.RespondError(err)
	}

	refreshClaims, err := parseToken(data.RefreshToken, "./public.pem")
	if err != nil {
		return c.RespondError(err)
	}

	// Check expiration times
	currentTime := time.Now().Unix()
	isAccessTokenExpired := accessClaims.ExpiresAt <= currentTime
	isRefreshTokenExpired := refreshClaims.ExpiresAt <= currentTime

	if isRefreshTokenExpired {
		return c.RespondError(fmt.Errorf("refresh token expired"))
	}

	if isAccessTokenExpired {
		return c.Respond(map[string]interface{}{
			"status": true,
			"data": map[string]string{
				"access_token":  data.AccessToken,
				"refresh_token": data.RefreshToken,
			},
			"error": nil,
		}, http.StatusOK)
	}

	// Compare user IDs
	if accessClaims.ID != refreshClaims.ID {
		return c.RespondError(fmt.Errorf("user ID mismatch between access and refresh tokens"))
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := commands.GenToken(user.AuthClaims{
		ID:   accessClaims.ID,
		Role: accessClaims.Role,
	}, "./private.pem")
	if err != nil {
		return c.RespondError(err)
	}

	// Respond with the new tokens
	return c.Respond(map[string]interface{}{
		"status": true,
		"data": map[string]string{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		},
		"error": nil,
	}, http.StatusOK)
}
func parseToken(tokenString, publicKeyPath string) (*user.AuthClaims, error) {
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &user.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*user.AuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
