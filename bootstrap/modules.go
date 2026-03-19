package bootstrap

import (
	"main/internal/application/http/v2"
	"main/internal/application/jobs"
	"main/internal/config"
	"main/pkg"

	chat_domain "main/internal/domain/chat"

	auth_infrastructure "main/internal/infrastructure/auth"
	chat_infrastructure "main/internal/infrastructure/chat"
	social_infrastructure "main/internal/infrastructure/social"

	auth_http "main/internal/application/http/v2/auth"
	chat_http "main/internal/application/http/v2/chat"
	health_http "main/internal/application/http/v2/health"
	swagger_http "main/internal/application/http/v2/swagger"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,

	chat_domain.Module,

	auth_infrastructure.Module,
	chat_infrastructure.Module,
	social_infrastructure.Module,

	http.Module,
	jobs.Module,
	auth_http.Module,
	health_http.Module,
	swagger_http.Module,
	chat_http.Module,
)
