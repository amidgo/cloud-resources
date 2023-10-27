package app

import (
	"github.com/amidgo/cloud-resources/config"
	"github.com/amidgo/cloud-resources/internal/executor"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/pricemanager"
	"github.com/amidgo/cloud-resources/internal/resourcemanager"
	"log"
)

func Executor(
	resourceType resourcetype.ResourceType,
	priceStorage pricemanager.PriceStorage,
	resourceStorage resourcemanager.ResourceStorage,
	statisticsStorage executor.StatiscticsStorage,
	deltaCounterFabric executor.DeltaCounterFabric,
	maxMachineCount int,
	cnf config.HealthLoadConfig,
) executor.Executor {
	// создаём priceManager для нашего типа ресурса
	priceManager := pricemanager.New(priceStorage, resourceType)
	// создаём resourceStorage который возвращает ресурсы заданного типа
	resourceStorage = &resourcemanager.ResourceStorageTypeWrapper{
		ResourceStorage: resourceStorage,
		Type:            resourceType,
	}
	// создаём resourceManager
	resourceManager := resourcemanager.NewResourceManager(resourceStorage, priceManager, maxMachineCount)

	// логгер которые добавляем префикс типа в каждый лог
	logger := executor.ResourceTypeLoggerWrapper{
		Logger:       log.Default(),
		ResourceType: resourceType,
	}
	// создаём все необходимые компоненты для создания executor
	initter := executor.NewInitter(resourceManager, resourceStorage, cnf.CPU.Init, cnf.RAM.Init)
	deltaExecutor := executor.NewDeltaExecutor(logger, resourceManager)
	deltaCounter := deltaCounterFabric.DeltaCounter(resourceType)

	// возвращаем executor
	return executor.NewExecutor(
		logger,
		initter,
		resourceStorage,
		statisticsStorage,
		deltaExecutor,
		deltaCounter,
	)
}
