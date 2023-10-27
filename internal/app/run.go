package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/amidgo/cloud-resources/config"
	"github.com/amidgo/cloud-resources/internal/model/resourcetype"
	"github.com/amidgo/cloud-resources/internal/scheduler"
	"github.com/amidgo/cloud-resources/internal/storage/http/pricestorage"
	"github.com/amidgo/cloud-resources/internal/storage/http/resourcestorage"
	"github.com/amidgo/cloud-resources/internal/storage/http/statisticsstorage"
	"github.com/amidgo/cloud-resources/pkg/httpclient"
)

func HttpClient(token string) *httpclient.HttpClient {

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	client := http.Client{
		Timeout:   time.Second * 5,
		Transport: transport,
	}

	httpClient := httpclient.NewHttpClient(
		&client,
		"https://mts-olimp-cloud.codenrock.com/api",
		func(r *http.Request) {
			q := r.URL.Query()
			q.Add("token", token)
			r.URL.RawQuery = q.Encode()
		},
	)
	return httpClient
}

func Run() {
	ctx := context.Background()

	cnf := config.ParseAppConfig()

	client := HttpClient(cnf.Api.Token)
	// изменяем output для наших логов
	// в нашем случае мы отсылаем логи на удалённый сервер
	// output := NewLogOutput(client.Client, cnf.Logger.LogURL)
	// log.SetOutput(output)

	statisticsStorage := statisticsstorage.New(client)
	resourceStorage := resourcestorage.New(client)
	priceStorage := pricestorage.New(client)

	fabric := DeltaCounterFabric(cnf)
	// создаём executor для DB
	dbExecutor := Executor(
		resourcetype.DB,
		priceStorage,
		resourceStorage,
		statisticsStorage,
		fabric,
		cnf.DB.MaxMachineCount,
		cnf.DB.HealthLoad,
	)
	// создаём executor для VM
	vmExecutor := Executor(
		resourcetype.VM,
		priceStorage,
		resourceStorage,
		statisticsStorage,
		fabric,
		cnf.VM.MaxMachineCount,
		cnf.VM.HealthLoad,
	)

	// создаём scheduler который будет вызывать два наших executorа раз в CheckTime
	scheduler := scheduler.New(time.Duration(cnf.Scheduler.CheckTime), dbExecutor, vmExecutor)

	log.Printf("service runned %s", time.Now().String())
	scheduler.Run(ctx)
}
