package billbee_importer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"billbee-importer/internal/billbee-importer/importer"
	"github.com/justtrackio/gosoline/pkg/cfg"
	"github.com/justtrackio/gosoline/pkg/http"
	"github.com/justtrackio/gosoline/pkg/kernel"
	"github.com/justtrackio/gosoline/pkg/log"
)

type StockxAliasBillbeeBridge struct {
	kernel.ForegroundModule
	logger                   log.Logger
	httpClient               http.Client
	billbeeEmailDelaySeconds int
	mailer                   Mailer
	importer                 []importer.Importer
}

func NewModule() func(ctx context.Context, config cfg.Config, logger log.Logger) (kernel.Module, error) {
	return func(ctx context.Context, config cfg.Config, logger log.Logger) (kernel.Module, error) {
		return NewStockxAliasBillbeeBridge(ctx, config, logger)
	}
}

func NewStockxAliasBillbeeBridge(ctx context.Context, config cfg.Config, logger log.Logger) (kernel.Module, error) {
	billbeeEmailDelaySeconds := config.GetInt("billbee_email_delay_seconds", 5)

	var runningImporter []importer.Importer
	for name, importerFactory := range importer.Sources {
		i, err := importerFactory(name, config, logger)

		if err != nil {
			println(fmt.Sprintf("failed to initialize %s importer - will skip", name))
			continue
		}

		runningImporter = append(runningImporter, i)
	}

	return &StockxAliasBillbeeBridge{
		logger:                   logger,
		httpClient:               http.NewHttpClient(config, logger),
		billbeeEmailDelaySeconds: billbeeEmailDelaySeconds,
		importer:                 runningImporter,
		mailer:                   NewMailer(config, logger),
	}, nil
}

func (e *StockxAliasBillbeeBridge) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := e.Start()

			if err != nil {
				println("critical error occurred")
			}

			println("Please enter the exit")
			input := bufio.NewScanner(os.Stdin)
			input.Scan()

			if err != nil {
				return err
			}

			println("Stockx/Alias sales to billbee import done")

			return nil
		}
	}
}

func (e *StockxAliasBillbeeBridge) Start() error {
	println("Billbee Importer provided by CaptainBarnius#0001")

	for _, i := range e.importer {
		sales := i.GetBillbeeSales()
		println(fmt.Sprintf("Start sending %d %s sales to billbee with a delay of %d seconds after each email to %s", len(sales), i.GetName(), e.billbeeEmailDelaySeconds, i.GetEmailEndpoint()))
		println("Please enter to start")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		for _, sale := range sales {
			err := e.SendEmail(sale, i.GetEmailEndpoint())

			if err == nil {
				println(fmt.Sprintf("%s - Email has been sent to %s", sale.OrderID, i.GetEmailEndpoint()))
			} else {
				println(fmt.Sprintf("%s - Failed to send email to %s", sale.OrderID, i.GetEmailEndpoint()))
			}
			time.Sleep(time.Duration(e.billbeeEmailDelaySeconds) * time.Second)
		}
	}

	return nil
}

func (e *StockxAliasBillbeeBridge) SendEmail(billbeeSale importer.BillbeeSale, emailEndpoint string) error {
	subject := fmt.Sprintf("Stockx Sale %s for Billbee", billbeeSale.OrderID)
	request := e.mailer.NewRequest(emailEndpoint, subject)

	err := request.Send(billbeeSale)

	return err
}
