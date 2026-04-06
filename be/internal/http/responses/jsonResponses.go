package responses

import "github.com/gofiber/fiber/v3"

const (
	ErrorMessageLocalKey = "response.error.message"
	ErrorDetailsLocalKey = "response.error.details"
)

// Error возвращает стандартный JSON-ответ с ошибкой.
//
// Параметры:
//
//	c — контекст Fiber.
//	status — HTTP-код из диапазона 400–599.
//	msg — сообщение ошибки.
//	details — необязательные детали ошибки.
//
// Возвращает:
//
//	Ошибку отправки HTTP-ответа.
func Error(c fiber.Ctx, status int, msg string, details ...any) error {
	if status < 400 || status >= 600 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "Некорректный код статуса для ошибки (должен быть 4xx или 5xx)",
		})
	}
	localizedMessage := localizeErrorMessage(msg)
	resp := fiber.Map{
		"ok":      false,
		"message": localizedMessage,
	}
	var detailValue any
	if len(details) > 0 && details[0] != nil {
		detailValue = details[0]
		resp["details"] = detailValue
	}
	c.Locals(ErrorMessageLocalKey, localizedMessage)
	if detailValue != nil {
		c.Locals(ErrorDetailsLocalKey, detailValue)
	}
	return c.Status(status).JSON(resp)
}

// Success возвращает стандартный JSON-ответ об успешной операции.
//
// Параметры:
//
//	c — контекст Fiber.
//	status — HTTP-код из диапазона 200–299.
//	data — полезная нагрузка ответа.
//
// Возвращает:
//
//	Ошибку отправки HTTP-ответа.
func Success(c fiber.Ctx, status int, data any) error {
	if status < 200 || status >= 300 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"message": "Некорректный код статуса для успешного ответа (должен быть 2xx)",
		})
	}
	return c.Status(status).JSON(fiber.Map{
		"ok":   true,
		"data": data,
	})
}
