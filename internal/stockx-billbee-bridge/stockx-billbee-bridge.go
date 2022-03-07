package stockx_billbee_bridge

import (
	"bufio"
	"context"
	"fmt"
	"github.com/applike/gosoline/pkg/cfg"
	"github.com/applike/gosoline/pkg/http"
	"github.com/applike/gosoline/pkg/kernel"
	"github.com/applike/gosoline/pkg/mon"
	"github.com/gocarina/gocsv"
	"os"
	"time"
)

type StockxBillbeeBridge struct {
	kernel.ForegroundModule
	logger                   mon.Logger
	httpClient           http.Client
	filePath             string
	billbeeEmailEndpoint string
	billbeeEmailDelaySeconds int
	stockxOrderStatusId      int
	vatRate                  int
	waitBeforeImportSeconds  int
	stockxPayoutMethod       string
	importComment            string
	internalNote             string
	stockxCurrency           string
	mailer                   Mailer
	stockxSales              []StockxSale
}

func New(filePath string) *StockxBillbeeBridge {
	return &StockxBillbeeBridge{
		filePath: filePath,
	}
}

func (e *StockxBillbeeBridge) Boot(config cfg.Config, logger mon.Logger) error {
	e.logger = logger
	e.httpClient = http.NewHttpClient(config, logger)
	e.billbeeEmailEndpoint = config.GetString("billbee_email_endpoint", "FILL_HERE")
	e.billbeeEmailDelaySeconds = config.GetInt("billbee_email_delay_seconds", 5)
	e.stockxPayoutMethod = config.GetString("stockx_payout_method", "paypal")
	e.stockxCurrency = config.GetString("stockx_currency", "EUR")
	e.importComment = config.GetString("optional_comment", "0")
	e.internalNote = config.GetString("internal_note", "Generated from Bridge Application")
	e.stockxOrderStatusId = config.GetInt("stockx_order_status_id", 7)
	e.waitBeforeImportSeconds = config.GetInt("wait_before_import_seconds", 60)
	e.vatRate = config.GetInt("vat_rate", 19)
	e.mailer = NewMailer(config, logger)

	if e.billbeeEmailEndpoint == "SHOP_ID-json-orderimport@inbound.billbee.de" {
		println("Please provide your Billbee Email Endpoint in the config.dist.yml file")

		println("Please enter the exit")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		return fmt.Errorf("billbee email endpoint is missing")
	}

	stockxSales, err := ParseCsvToStruct(e.filePath)

	if err != nil {
		println("Error occurred while parsing stockx.csv file. Please check it")

		println("Please enter the exit")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()


		return err
	}

	e.stockxSales = stockxSales

	return nil
}

func ParseCsvToStruct(filePath string) ([]StockxSale, error) {
	salesFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer salesFile.Close()

	var stockxSales []StockxSale
	err = gocsv.UnmarshalFile(salesFile, &stockxSales)

	if err != nil {
		return nil, err
	}

	return stockxSales, nil
}

func (e *StockxBillbeeBridge) Run(ctx context.Context) error {
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

			println("Stockx sales to billbee import done")

			return nil
		}
	}
}

func (e *StockxBillbeeBridge) Start() error {
	println("Stockx Billbee Bridge provided by CaptainBarnius#0001")
	println(fmt.Sprintf("Start sending %d stockx sales to billbee with a delay of %d seconds after each email to %s", len(e.stockxSales), e.billbeeEmailDelaySeconds, e.billbeeEmailEndpoint))
	println(fmt.Sprintf("You have now %d seconds to abort the import by terminating this application.", e.waitBeforeImportSeconds))

	time.Sleep(time.Duration(e.waitBeforeImportSeconds) * time.Second)

	for _, stockxSale := range e.stockxSales {
		orderItems := []BillbeeOrderItem{
			{
				Quantity:   stockxSale.Quantity,
				Totalprice: stockxSale.ListingPrice,
				ProductID:  nil,
				Name:       stockxSale.SkuName,
				Sku:        "",
				Imageurl:   "",
				Taxrate:    e.vatRate,
			},
			{
				Quantity:   1,
				Totalprice: negate(stockxSale.SellerFee),
				ProductID:  nil,
				Name:       "./. Transaktionsgebühr",
				Sku:        "",
				Imageurl:   "",
				Taxrate:    e.vatRate,
			},
			{
				Quantity:   1,
				Totalprice: negate(stockxSale.PaymentProcessing),
				ProductID:  nil,
				Name:       "./. Zahlungsabwicklungsgebühr",
				Sku:        "",
				Imageurl:   "",
				Taxrate:    e.vatRate,
			},
		}

		billbeeSale := BillbeeSale{
			OrderID:       stockxSale.OrderNumber,
			CurrencyCode:  e.stockxCurrency,
			PaymentMethod: e.stockxPayoutMethod,
			OrderDate:     stockxSale.SaleDate,
			UstID:         stockxSale.SoldToVatNumber,
			Email:         "",
			Telephone:     "",
			Invoice: BillbeeInvoice{
				Salutation:  0,
				Firstname:   "StockX LLC",
				Lastname:    "",
				Housenumber: "1046",
				Address1:    "Woodward Avenue",
				Address2:    "",
				Zip:         "48226",
				State:       "MI",
				City:        "Detroit",
				Country:     "UNITED STATES OF AMERICA",
				Company:     "StockX LLC",
			},
			Shipping: BillbeeShipping{
				Salutation:  0,
				Firstname:   "StockX LLC",
				Lastname:    "",
				Housenumber: "",
				Address1:    "DE RUN 4312",
				Address2:    "SMOT",
				Zip:         "5503LN",
				State:       "",
				City:        "VELDHOVEN",
				Country:     "NETHERLANDS",
				Company:     "",
			},
			OrderStatusID: e.stockxOrderStatusId,
			InternalNote:  e.internalNote,
			Comment:       e.importComment,
			Shipcost:      0,
			ShippingDate:  stockxSale.PayoutDate,
			OrderItems:    orderItems,
			Tags:          nil,
		}

		err := e.SendEmail(billbeeSale)

		if err == nil {
			println(fmt.Sprintf("%s - Email has been sent to %s", billbeeSale.OrderID, e.billbeeEmailEndpoint))
		} else {
			println(fmt.Sprintf("%s - Failed to send email to %s", billbeeSale.OrderID, e.billbeeEmailEndpoint))
		}

		time.Sleep(time.Duration(e.billbeeEmailDelaySeconds) * time.Second)
	}

	return nil
}

func (e *StockxBillbeeBridge) SendEmail(billbeeSale BillbeeSale) error {
	subject := fmt.Sprintf("Stockx Sale %s for Billbee", billbeeSale.OrderID)
	request := e.mailer.NewRequest(e.billbeeEmailEndpoint, subject)

	err := request.Send(billbeeSale)

	return err
}

func negate(val float64) float64 {
	return val * -1
}
