package pkg

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRequestHandler),
	fx.Provide(GetLogger),
	fx.Provide(NewGormDatabase),
	fx.Provide(NewNatsClient),
	fx.Provide(NewGrpcClient),
)
