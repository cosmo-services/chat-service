package bootstrap

import (
	api "main/internal/application/http/v2"
	"main/internal/application/jobs"
	"main/internal/config"
	"main/pkg"

	chat_domain "main/internal/domain/chat"

	auth_infrastructure "main/internal/infrastructure/auth"
	chat_infrastructure "main/internal/infrastructure/chat"

	health_http "main/internal/application/http/v2/health"
	swagger_http "main/internal/application/http/v2/swagger"
	test_http "main/internal/application/http/v2/test"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,

	chat_domain.Module,

	auth_infrastructure.Module,
	chat_infrastructure.Module,

	api.Module,
	jobs.Module,
	health_http.Module,
	swagger_http.Module,
	test_http.Module,
)
