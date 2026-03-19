package http

import (
	chat_api "main/internal/application/http/v2/chat"
	health_api "main/internal/application/http/v2/health"
	swagger_api "main/internal/application/http/v2/swagger"

	"go.uber.org/fx"
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	healthRoutes *health_api.HealthRoutes,
	swaggerRoutes *swagger_api.SwaggerRoutes,
	chatRoutes *chat_api.ChatRoutes,
) Routes {
	return Routes{
		healthRoutes,
		swaggerRoutes,
		chatRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}

var Module = fx.Options(
	fx.Provide(NewRoutes),
)
