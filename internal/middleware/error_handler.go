package middleware

import (
	"github.com/labstack/echo/v4"
	appErrors "github.com/omohide_map_backend/pkg/errors"
)

func CustomErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	appErr := appErrors.GetAppError(err)

	response := map[string]interface{}{
		"error": map[string]interface{}{
			"message": appErr.Message,
		},
	}

	if appErr.Detail != "" {
		response["error"].(map[string]interface{})["detail"] = appErr.Detail
	}

	c.JSON(appErr.HTTPStatus(), response)
}
