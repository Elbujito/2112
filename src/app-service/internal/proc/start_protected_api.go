package proc

import (
	"github.com/Elbujito/2112/internal/api/routers"
	"github.com/Elbujito/2112/internal/clients/service"
)

func StartProtectedApi() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	routers.InitProtectedAPIRouter()
	routers.ProtectedAPIRouter().Start(config.Host, config.ProtectedApiPort)
}