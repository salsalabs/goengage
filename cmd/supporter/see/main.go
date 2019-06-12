package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app   = kingpin.New("see-supporter", "A command-line app to to show supporters for an email.")
		login = app.Flag("login", "YAML file with API token").Required().String()
		email = app.Flag("email", "Email address to look up").Required().String()
	)
	app.Parse(os.Args[1:])
	if login == nil || len(*login) == 0 {
		fmt.Println("Error --login is required.")
		os.Exit(1)
	}
	if email == nil || len(*email) == 0 {
		fmt.Println("Error --email is required.")
		os.Exit(1)
	}
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}

	a := []string{*email}
	rqt := goengage.SupSearchIDRequest{
		Identifiers:    a,
		IdentifierType: "EMAIL_ADDRESS",
		Offset:         0,
		Count:          int32(len(a)),
	}
	var resp goengage.SupSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SupSearchMethod,
		Fragment: goengage.SupSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err = n.Do()
	if err != nil {
		panic(err)
	}
	for _, s := range resp.Payload.Supporters {
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
}
