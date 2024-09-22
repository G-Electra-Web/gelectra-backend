package handlers

import (
	"github.com/gelectra/gelectra-backend/internal/database"
	"github.com/gofiber/fiber/v2"
)

func ListProjects(c *fiber.Ctx) error {
	var approvedProjects []database.Project

	// Fetch all projects with status 'approved'
	if err := database.DB.Where("status = ?", "approved").Preload("Users").Find(&approvedProjects).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve approved projects"})
	}

	// If no approved projects found
	if len(approvedProjects) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No approved projects found"})
	}

	return c.JSON(fiber.Map{
		"approved_projects": approvedProjects,
	})
}

func GetProject(c *fiber.Ctx) error {
	projectID := c.Params("id")
	var project database.Project

	if err := database.DB.Table("projects").First(&project, projectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	return c.JSON(project)
}

func GetHomePageProjects(c *fiber.Ctx) error {
	var projects []database.Project

	// Fetch projects set as homepage projects
	if err := database.DB.Where("is_home_page_project = ?", true).Find(&projects).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve homepage projects",
		})
	}

	return c.JSON(fiber.Map{
		"home_page_projects": projects,
	})
}

func GetSortedMembers(c *fiber.Ctx) error {
	var members []struct {
		FullName string `json:"full_name"`
		Role     string `json:"role"`
	}

	// Fetch members sorted by the Role field, selecting only FullName and Role
	if err := database.DB.Table("users").Select("full_name, role").Order("role ASC").Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve members",
		})
	}

	return c.JSON(fiber.Map{
		"members": members,
	})
}
