package main

import (
	billbee_importer "billbee-importer/internal/billbee-importer"
	"github.com/justtrackio/gosoline/pkg/application"
)

func main() {
	application.Run(
		application.WithModuleFactory("billbee-importer", billbee_importer.NewModule()),
	)
}
