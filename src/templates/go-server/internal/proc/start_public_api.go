package proc

import (
	"github.com/Elbujito/2112/src/templates/go-server/internal/api/routers"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/service"
	"github.com/Elbujito/2112/src/templates/go-server/internal/config"
)

func StartPublicApi() {
	serviceCli := service.GetClient()
	c := serviceCli.GetConfig()
	publicApiRouter := routers.InitPublicAPIRouter(config.Env)
	publicApiRouter.Start(c.Host, c.PublicApiPort)
}
