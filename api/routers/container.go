package routers

import (
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/decentorganization/topaz/api/services"
	"github.com/decentorganization/topaz/shared/database"
)

// ProvideAppController ...
func ProvideAppController() controllers.AppController {
	return controllers.NewAppController(ProvideAppService())
}

// ProvideAppService ...
func ProvideAppService() services.AppService {
	return services.NewAppService(ProvideDataManager())
}

var dataManager = database.NewDataManager(database.Manager)

// ProvideDataManager ...
func ProvideDataManager() database.DataManager {
	return dataManager
}
