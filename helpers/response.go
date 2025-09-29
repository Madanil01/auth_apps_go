package	helpers

import (
	"github.com/gofiber/fiber/v2"
	// "apps_v1/models"
)

func Response(c *fiber.Ctx, code int, message string, data interface{}) error {
    statusMap := map[int]int{
        200: fiber.StatusOK,
        201: fiber.StatusCreated,
        400: fiber.StatusBadRequest,
        401: fiber.StatusUnauthorized,
        403: fiber.StatusForbidden,
        404: fiber.StatusNotFound,
        500: fiber.StatusInternalServerError,
    }

    // default status OK
    status, ok := statusMap[code]
    if !ok {
        status = fiber.StatusOK
    }

    return c.Status(status).JSON(fiber.Map{
        "status":  status,
        "message": message,
        "payload":  data,
    })
}
