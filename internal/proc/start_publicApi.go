package proc

import (
	"github.com/Elbujito/2112/internal/api/routers"
	"github.com/Elbujito/2112/internal/clients/service"
)

func StartPublicApi() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	routers.InitPublicAPIRouter()
	routers.PublicAPIRouter().Start(config.Host, config.PublicApiPort)
}
