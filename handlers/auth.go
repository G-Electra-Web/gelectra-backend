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

func InitAdmin() error {
	adminEmail := k.String("admin.email")
	adminPassword := k.String("admin.password")
	adminFullName := k.String("admin.fullname")

	var user database.User
	if err := database.DB.Where("email = ?", adminEmail).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User does not exist, create it
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			newUser := database.User{
				Email:        adminEmail,
				PasswordHash: string(hashedPassword),
				FullName:     adminFullName,
				Role:         "admin",
				JoinDate:     time.Now(),
			}

			if err := database.DB.Create(&newUser).Error; err != nil {
				return err
			}

			// Create associated Admin record
			newAdmin := database.Admin{
				UserID:     newUser.ID,
				Privileges: "all", // You might want to adjust this based on your requirements
			}

			if err := database.DB.Create(&newAdmin).Error; err != nil {
				return err
			}

			return nil
		}
		return err
	}

	// User exists, check if Admin record exists
	var admin database.Admin
	if err := database.DB.Where("user_id = ?", user.ID).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Admin record doesn't exist, create it
			newAdmin := database.Admin{
				UserID:     user.ID,
				Privileges: "all", // You might want to adjust this based on your requirements
			}

			if err := database.DB.Create(&newAdmin).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find user in the database
	var user database.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := k.String("auth.jwt_secret")
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// If the user is an admin, fetch admin details
	var adminInfo *database.Admin
	if user.Role == "admin" {
		var admin database.Admin
		if err := database.DB.Where("user_id = ?", user.ID).First(&admin).Error; err == nil {
			adminInfo = &admin
		}
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":         user.ID,
			"email":      user.Email,
			"full_name":  user.FullName,
			"role":       user.Role,
			"join_date":  user.JoinDate,
			"admin_info": adminInfo,
		},
	})
}

func SignUp(c *fiber.Ctx) error {
	type SignUpInput struct {
		Password        string `json:"password" validate:"required"`
		FullName        string `json:"full_name" validate:"required"`
		Email           string `json:"email" validate:"required,email"`
		RegistrationNum string `json:"registration_num"`
	}

	var input SignUpInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the email is whitelisted
	var whitelistedUser database.User
	if err := database.DB.Where("email = ?", input.Email).First(&whitelistedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email not whitelisted"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Check if the user already exists and has completed registration
	var existingUser database.User
	if err := database.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		if existingUser.PasswordHash != "" || existingUser.FullName != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User already registered"})
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Update the user
	updatedUser := database.User{
		ID:           whitelistedUser.ID,
		PasswordHash: string(hashedPassword),
		Email:        input.Email,
		FullName:     input.FullName,
		Role:         whitelistedUser.Role,
		JoinDate:     time.Now(),
	}

	if err := tx.Save(&updatedUser).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// If the user is a member, create a Member record
	if updatedUser.Role == "member" {
		if input.RegistrationNum == "" {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Registration number is required for members"})
		}

		newMember := database.Member{
			UserID:          updatedUser.ID,
			RegistrationNum: input.RegistrationNum,
			MembershipLevel: "regular", // You might want to adjust this based on your requirements
			JoinDate:        time.Now(),
		}

		if err := tx.Create(&newMember).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create member record"})
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User signed up successfully",
		"user": fiber.Map{
			"id":        updatedUser.ID,
			"email":     updatedUser.Email,
			"full_name": updatedUser.FullName,
			"role":      updatedUser.Role,
			"join_date": updatedUser.JoinDate,
		},
	})
}
