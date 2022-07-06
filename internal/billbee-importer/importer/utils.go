package importer

import "strings"

func negate(val float64) float64 {
	return val * -1
}

// convert tt/mm/yyyy to yyyy-mm-tt
func convertAliasDate(date string) string {
	parts := strings.Split(date, "/")

	return "20" + parts[2] + "-" + parts[1] + "-" + parts[0]
}
