package bootstrap

import (
	api "main/internal/application/http/v2"
	"main/internal/application/jobs"
	"main/internal/config"
	"main/pkg"

	auth_infrastructure "main/internal/infrastructure/auth"
	test_infrastructure "main/internal/infrastructure/test"

	health_http "main/internal/application/http/v2/health"
	swagger_http "main/internal/application/http/v2/swagger"
	test_http "main/internal/application/http/v2/test"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	config.Module,
	pkg.Module,

	auth_infrastructure.Module,
	test_infrastructure.Module,

	api.Module,
	jobs.Module,
	health_http.Module,
	swagger_http.Module,
	test_http.Module,
)
