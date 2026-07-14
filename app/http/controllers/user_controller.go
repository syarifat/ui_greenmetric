package controllers

import (
	"ui_greenmetric/app/facades"
	"ui_greenmetric/app/models"

	"github.com/goravel/framework/contracts/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// Index lists users belonging to the admin's campus (or all for SUPER_ADMIN)
func (r *UserController) Index(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	var users []models.User
	query := facades.Orm().Query()

	if currentUser.Role == "ADMIN_KAMPUS" {
		query = query.Where("campus_id = ?", currentUser.CampusID)
	}

	err := query.Get(&users)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to fetch users",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

// Store creates a new operator user
func (r *UserController) Store(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	role := ctx.Request().Input("role")

	if name == "" || email == "" || password == "" || role == "" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Name, Email, Password, and Role are required",
		})
	}

	// Validate role values
	validRoles := map[string]bool{
		"ADMIN_KAMPUS": true, "OPERATOR_SI": true, "OPERATOR_EC": true,
		"OPERATOR_WS": true, "OPERATOR_WR": true, "OPERATOR_TR": true,
		"OPERATOR_ED": true, "OPERATOR_GD": true,
	}
	if !validRoles[role] {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid role specified",
		})
	}

	// Determine campus_id
	var campusID uint
	if currentUser.Role == "SUPER_ADMIN" {
		campusID = uint(ctx.Request().InputInt("campus_id"))
		if campusID == 0 {
			return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
				"status":  "error",
				"code":    http.StatusUnprocessableEntity,
				"message": "campus_id is required for SUPER_ADMIN",
			})
		}
	} else {
		campusID = currentUser.CampusID
	}

	// Check if email already exists
	count, err := facades.Orm().Query().Model(&models.User{}).Where("email = ?", email).Count()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Database error",
		})
	}
	if count > 0 {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "Email already registered",
		})
	}

	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Hashing error",
		})
	}

	user := models.User{
		CampusID: campusID,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to create user",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"status":  "success",
		"message": "User created successfully",
		"data":    user,
	})
}

// Update modifies an existing user
func (r *UserController) Update(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	id := ctx.Request().RouteInt("id")
	var user models.User
	err := facades.Orm().Query().Where("id = ?", id).First(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
	}

	// Restrict multi-tenant checks
	if currentUser.Role == "ADMIN_KAMPUS" && user.CampusID != currentUser.CampusID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: You cannot modify users from another campus",
		})
	}

	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	role := ctx.Request().Input("role")

	if name != "" {
		user.Name = name
	}
	if email != "" {
		// Check unique email
		count, err := facades.Orm().Query().Model(&models.User{}).Where("email = ? AND id != ?", email, id).Count()
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Database error",
			})
		}
		if count > 0 {
			return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
				"status":  "error",
				"code":    http.StatusUnprocessableEntity,
				"message": "Email already registered by another user",
			})
		}
		user.Email = email
	}

	if password != "" {
		hashed, err := facades.Hash().Make(password)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"status":  "error",
				"code":    http.StatusInternalServerError,
				"message": "Hashing error",
			})
		}
		user.Password = hashed
	}

	if role != "" {
		validRoles := map[string]bool{
			"ADMIN_KAMPUS": true, "OPERATOR_SI": true, "OPERATOR_EC": true,
			"OPERATOR_WS": true, "OPERATOR_WR": true, "OPERATOR_TR": true,
			"OPERATOR_ED": true, "OPERATOR_GD": true,
		}
		if validRoles[role] {
			user.Role = role
		}
	}

	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to update user",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "User updated successfully",
		"data":    user,
	})
}

// Destroy deletes an existing user
func (r *UserController) Destroy(ctx http.Context) http.Response {
	var currentUser models.User
	if err := facades.Auth(ctx).User(&currentUser); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"status":  "error",
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
	}

	id := ctx.Request().RouteInt("id")
	var user models.User
	err := facades.Orm().Query().Where("id = ?", id).First(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"status":  "error",
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
	}

	// Restrict multi-tenant checks
	if currentUser.Role == "ADMIN_KAMPUS" && user.CampusID != currentUser.CampusID {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"status":  "error",
			"code":    http.StatusForbidden,
			"message": "Forbidden: You cannot delete users from another campus",
		})
	}

	// Prevent user deleting themselves
	if currentUser.ID == user.ID {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"status":  "error",
			"code":    http.StatusUnprocessableEntity,
			"message": "You cannot delete your own account",
		})
	}

	if _, err := facades.Orm().Query().Delete(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"status":  "error",
			"code":    http.StatusInternalServerError,
			"message": "Failed to delete user",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
