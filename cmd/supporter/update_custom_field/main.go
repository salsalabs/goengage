package main

import (
	"encoding/json"
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

//Program entry point.  Look for supporters with an email.  Errors are noisy and fatal.
func main() {
	var (
		app        = kingpin.New("see-supporter", "A command-line app to modify a custom field.")
		login      = app.Flag("login", "YAML file with API token").Required().String()
		email      = app.Flag("email", "Supporter's email address").Required().String()
		fieldName  = app.Flag("fieldName", "Custom field name to modify").Required().String()
		fieldValue = app.Flag("fieldValue", "Value to assign to `fieldName`").Required().String()
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
	s, err := goengage.SupporterByEmail(e, *email)
	if err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println("--------------- Supporter Found ----------------")
	fmt.Print(string(b))
	fmt.Println("")

	fmt.Printf("Changing %s to '%s'\n", *fieldName, *fieldValue)
	found := false
	for _, c := range s.CustomFieldValues {
		if c.Name == *fieldName {
			c.Value = *fieldValue
			found = true
		}
	}
	if !found {
		x := fmt.Sprintf("error: '%v' is not a valid custom field name", *fieldName)
		fmt.Println(x)
		fmt.Println("Please choose from one of these field names")
		for _, c := range s.CustomFieldValues {
			fmt.Printf("* %s\n", c.Name)
		}
		return
	}
	result, err := goengage.SupporterUpsert(e, s)
	if err != nil {
		fmt.Printf("Upsert failed with %s\n", err)
	}
	b, err = json.MarshalIndent(result, "", "    ")
	if err != nil {
		fmt.Printf("JSON marshall error, %s\n", err)
		fmt.Printf("Supporter si %+v\n", result)
	} else {
		fmt.Println("--------------- Supporter Results ----------------")
		fmt.Print(string(b))
		fmt.Println("")
	}
}
