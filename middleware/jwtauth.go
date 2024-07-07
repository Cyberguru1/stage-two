package middleware


import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/cyberguru1/stage-two/utils"
	"github.com/cyberguru1/stage-two/config"
)

func IsAuthorize(config *config.Config) func(ctx *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(config.Jwt.Secret),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "Bad request",
				"message": "Authentication failed",
				"statusCode": 401,
			})

			return nil
		},
	})

}

func GetUserIdFromContext(ctx *fiber.Ctx) (string, error) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	strId := claims["sub"].(string)

	return strId, nil
}

func ClaimToken(id uuid.UUID) (string, error) {
	config := config.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 1).Unix() // make it valid for a day

	s, err := token.SignedString([]byte(config.Jwt.Secret))

	if err != nil {
		utils.Errorf("error: ", err)
		return "", err
	}

	return s, nil

}
