package main

//Application scan for fundraising activities with dedications
//and write them to a CSV.
import (
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app   = kingpin.New("dedications", "Write dedications to a CSV")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	var service goengage.Service
	err = goengage.ReportFundraising(e, service)
	if err != nil {
		panic(err)
	}
}
