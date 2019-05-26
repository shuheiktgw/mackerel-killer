// +build integration

package mkk

import "os"

var (
	integrationMackerelToken = os.Getenv("MACKEREL_API_TOKEN")
	integrationMkk           *Mkk
)

func init() {
	integrationMkk = NewMkk(integrationMackerelToken)
}
