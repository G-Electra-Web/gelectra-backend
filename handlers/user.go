package handlers

import (
	"time"

	"github.com/gelectra/gelectra-backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateProject(c *fiber.Ctx) error {
	type CreateProjectInput struct {
		LeaderID    int    `json:"leader_id" validate:"required"`
		ProjectName string `json:"project_name" validate:"required"`
		Description string `json:"description"`
		ExtraLink   string `json:"extra_link"`
		UserIDs     []int  `json:"user_ids"` // List of user IDs to add to the project
	}

	var input CreateProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Generate a unique join code
	joinCode := uuid.New().String()

	// Create the project
	project := database.Project{
		LeaderID:    input.LeaderID,
		ProjectName: input.ProjectName,
		Description: input.Description,
		ExtraLink:   input.ExtraLink,
		Status:      "pending",
		CreatedBy:   input.LeaderID,
		JoinCode:    joinCode,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create project"})
	}

	// Add users to the project if provided
	if len(input.UserIDs) > 0 {
		var users []database.User
		if err := database.DB.Where("user_id IN ?", input.UserIDs).Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add users to the project"})
		}
		if err := database.DB.Model(&project).Association("Users").Append(&users); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to associate users with the project"})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(project)
}

func JoinProject(c *fiber.Ctx) error {
	type JoinProjectInput struct {
		UserID   int    `json:"user_id" validate:"required"`
		JoinCode string `json:"join_code" validate:"required"`
	}

	var input JoinProjectInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var project database.Project
	if err := database.DB.Where("join_code = ?", input.JoinCode).First(&project).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid join code"})
	}

	// Check if user already joined
	var user database.User
	if err := database.DB.First(&user, input.UserID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := database.DB.Model(&project).Association("Users").Append(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to join the project"})
	}

	return c.JSON(fiber.Map{
		"message": "User successfully joined the project",
	})
}

func EditProfile(c *fiber.Ctx) error {
	// Assuming user ID is retrieved from session or JWT token, otherwise provide it in the request
	userID, err := c.ParamsInt("user_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Define a struct to capture input fields (excluding email)
	type UpdateProfileInput struct {
		Username string `json:"username" validate:"required"`
		FullName string `json:"full_name" validate:"required"`
		Password string `json:"password,omitempty"` // Password is optional
		Facebook string `json:"facebook"`           // New field for Facebook link
		Twitter  string `json:"twitter"`            // New field for Twitter link
		LinkedIn string `json:"linkedin"`           // New field for LinkedIn link
	}

	var input UpdateProfileInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the user
	var user database.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Update the user's profile (excluding email)
	user.Username = input.Username
	user.FullName = input.FullName
	user.Facebook = input.Facebook
	user.Twitter = input.Twitter
	user.LinkedIn = input.LinkedIn

	// If the password is provided, update it (you should hash the password before saving)
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Save the updated user profile
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update profile"})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
