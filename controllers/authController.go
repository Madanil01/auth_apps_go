package controllers

import (
	"apps_v1/database"
	"apps_v1/models"
	"apps_v1/utils"
	"encoding/base64"
	"log"
	"os"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func Hello(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func Register(c *fiber.Ctx) error {
	var data map[string]string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err :=c.BodyParser(&data); err != nil {
		return err
	}
	var userCek models.User
	result := database.DB.Where("email = ?", data["email"]).First(&userCek)
	if result.Error == nil {
		// Jika user ditemukan, berarti email sudah ada
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"status":  fiber.StatusConflict,
			"message": "Email already exists",
		})
	} else{
	
		if data["password"] != data["confirm_pass"] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Password and confirm password do not match",
			})
		}
		password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Error hashing password",
			})
		}
		user := models.User{
			Username: data["username"],
			Name: data["name"],
			Email: data["email"],
			Password: password,
			Alamat: data["alamat"],
		}
		var email_encode string= base64.StdEncoding.EncodeToString([]byte(user.Email))
		var BASE_HOST string = os.Getenv("BASE_HOST")
		var url string = BASE_HOST + "/api/confirm/" + email_encode
		utils.SendEmail(data["email"], "Selamat bergabung di aplikasi", "Selamat bergabung di aplikasi "+url)
		database.DB.Create(&user)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":  fiber.StatusCreated,
			"message": "User registered successfully",
			"data": fiber.Map{
				"email": user.Email,
				"name": user.Name,
			},
		})
	}
	
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	// Durasi token: 1 jam (3600 detik)
	accessToken, refreshToken, err := utils.GenerateTokens(user.Id, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
		})
	}

	// Set refresh token ke cookie (biasanya HttpOnly)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	// Optionally: juga bisa set access token via cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
	})
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"status":  fiber.StatusOK,
		"payload": fiber.Map{
			"user": fiber.Map{
				"name":          user.Name,
				"email":         user.Email,
				"confirm_email": user.ConfirmEmail,
				"is_admin":      user.IsAdmin,
			},
			"access_token":  fiber.Map{
				"value":    accessToken,
				"expires":  time.Now().Add(120 * time.Minute),
				"HTTPOnly": true,
				"Secure":   false,
				"SameSite": "Lax",
			},
			"refresh_token": fiber.Map{
				"value":    refreshToken,
				"expires":  time.Now().Add(7 * 24 * time.Hour),
				"HTTPOnly": true,
				"Secure":   false,
				"SameSite": "Lax",
			},
		},
	})
	// return c.JSON(fiber.Map{
	// 	"message": "Login successful",
	// 	"status":  fiber.StatusOK,
	// 	"payload": fiber.Map{
	// 		"user": fiber.Map{
	// 		"name":          user.Name,
	// 		"email":         user.Email,
	// 		"confirm_email": user.ConfirmEmail,
	// 		"is_admin":      user.IsAdmin,
	// 	},
	// 	"token": fiber.Map{"value": accessToken,
	// 			"expires":  time.Now().Add(tokenDuration),
	// 			"HTTPOnly": true,
	// 			"Secure":   false, 
	// 			"SameSite": "Lax",
	// 				}, 
	// 	},
		
	// })
}

func ConfirmEmail(c *fiber.Ctx) error {
	// Ambil email dari URL
	email := c.Params("email")
	emailBytes, _ := base64.URLEncoding.DecodeString(email)
	emailStr := string(emailBytes) 
	// Cek apakah email kosong (tidak valid)
	if emailStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email parameter is required",
		})
	}
	// Cari user berdasarkan email
	var user models.User
	result := database.DB.Where("email = ?", emailStr).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Update status email terkonfirmasi
	user.ConfirmEmail = true
	database.DB.Save(&user)

	return c.JSON(fiber.Map{
		"message": "Email confirmed successfully",
		"user": fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
func Profile(c *fiber.Ctx) error {
	// Contoh: Ambil user_id dari JWT token
	userID := c.Params("user_id")

	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User profile",
		"user": fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
func Logout(c *fiber.Ctx) error {
	// Logout hanya mengembalikan respons sukses karena JWT tidak bisa langsung dihapus dari server
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	// Contoh: Ambil user_id dari JWT token
	userID := c.Params("user_id")

	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Update data user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Save(&user)

	return c.JSON(fiber.Map{
		"message": "User profile updated",
		"user": fiber.Map{
			"id":    user.Id,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing refresh token",
		})
	}
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return utils.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	email := claims["email"].(string)

	// Buat access token baru
	accessToken, _, err := utils.GenerateTokens(userID, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(1 * time.Minute),
		HTTPOnly: false,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"access_token": fiber.Map{
			"value":  accessToken,
			"expire": time.Now().Add(1 * time.Minute),
		},
	})
}
