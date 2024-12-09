package proc

import (
	"github.com/Elbujito/2112/src/app-service/internal/api/routers"
	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
)

func StartHiddenApi() {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	routers.InitHiddenAPIRouter()
	routers.HiddenAPIRouter().Start(config.Host, config.HiddenApiPort)
}
