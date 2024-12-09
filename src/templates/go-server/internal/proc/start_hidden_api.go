package proc

import (
	"github.com/Elbujito/2112/src/templates/go-server/internal/api/routers"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/service"
)

func StartHiddenApi() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	routers.InitHiddenAPIRouter()
	routers.HiddenAPIRouter().Start(config.Host, config.HiddenApiPort)
}
