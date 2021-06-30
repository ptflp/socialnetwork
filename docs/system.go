package docs

import (
	"gitlab.com/InfoBlogFriends/server/config"
)

// swagger:route GET /system/config system systemGetConfig
// Получение конфигурации сервера.
// security:
//   - Bearer: []
// responses:
//   200: systemGetConfig

// swagger:response systemGetConfig
type systemGetConfig struct {
	// in:body
	Body config.Config
}
