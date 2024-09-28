package handlers

import (
	"os"
	"time"

	"github.com/gelectra/gelectra-backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var k = koanf.New(".")

func init() {
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		panic(err)
	}
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the input username and password match the admin credentials from config.toml
	adminUsername := k.String("admin.username")
	adminPassword := k.String("admin.password")

	if input.Username == adminUsername && input.Password == adminPassword {
		// Generate JWT token for admin
		claims := jwt.MapClaims{
			"user_id": "admin",
			"role":    "admin",
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		return c.JSON(fiber.Map{"token": tokenString})
	}

	// Find user in the database
	var user database.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare the password (assuming you have hashed it)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token for regular user
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString, "user": user.Role})
}

func SignUp(c *fiber.Ctx) error {
	type SignUpInput struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		FullName string `json:"full_name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
	}

	var input SignUpInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the email is whitelisted
	var whitelistedUser database.User
	if err := database.DB.Table("users").Where("email = ?", input.Email).First(&whitelistedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email not whitelisted"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if the user already exists
	var existingUser database.User
	if err := database.DB.Table("users").Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		// Check if other columns are empty
		if existingUser.Username != "" || existingUser.PasswordHash != "" || existingUser.FullName != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Other fields must be empty"})
		}
	}

	// Create the new user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	newUser := database.User{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
		Email:        input.Email,
		FullName:     input.FullName,
		Role:         whitelistedUser.Role, // Take the role from the whitelisted user
		JoinDate:     time.Now(),
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User signed up successfully",
		"user":    newUser,
	})
}
