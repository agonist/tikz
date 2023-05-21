package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/agonist/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	userID, err := ValidateToken(token)

	if err != nil {
		fmt.Println("failed to parse token:", err)
		return err
	}

	fmt.Println(userID)

	c.Context().SetUserValue("UserID", userID)

	return c.Next()
}

type JWTClaim struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

const JWT_EXPIRATION = 24 * 7 * time.Hour

func GenerateJWT(user *types.User) (string, error) {
	var claims = JWTClaim{
		user.ID,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWT_EXPIRATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(signedToken string) (uint, error) {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) { return []byte(os.Getenv("JWT_SECRET")), nil })
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return 0, fmt.Errorf("unauthorized")
	}

	return claims.ID, nil
}
