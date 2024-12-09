package proc

import (
	"github.com/Elbujito/2112/src/templates/go-server/internal/api/routers"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/service"
)

func StartProtectedApi() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	routers.InitProtectedAPIRouter()
	routers.ProtectedAPIRouter().Start(config.Host, config.ProtectedApiPort)
}
