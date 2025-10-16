package middlewares

import (
	"log"
	"pleiades/server/models"
	"strings"

	auth "pleiades/gen/auth"
	"pleiades/server/config/grpc"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func IsAuthenticated(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Apiresponse{
			Status:  false,
			Message: "Authorization header is missing.",
		})
	}

	tokenString := strings.TrimSpace(authHeader)
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Apiresponse{
			Status:  false,
			Message: "Bearer token missing or invalid format.",
		})
	}

	md := metadata.New(map[string]string{
		"authorization": authHeader,
	})
	ctx := metadata.NewOutgoingContext(c.Context(), md)

	resp, err := grpc.AuthService.IsAuthenticated(ctx, &auth.AuthRequest{
		RequestId: uuid.New().String(),
	})
	if err != nil {
		log.Printf("❌ AuthService.IsAuthenticated error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.Apiresponse{
			Status:  false,
			Message: "Invalid or expired token.",
		})
	}

	c.Locals("user_id", resp.UserId)
	c.Locals("role", resp.Role)
	c.Locals("is_verified", resp.IsVerified)

	return c.Next()
}
