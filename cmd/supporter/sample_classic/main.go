package main

import "github.com/salsalabs/goengage"

//App to read a number of supporter records from Salsa and
//write them to Engage.

func xform(c map[string]string) *goengage.Supporter {
	// I can't find a place in engage to store job-related info.
	// leaving it out of this test.

	s := goengage.Supporter{
		FirstName:        c["First_Name"],
		LanguageCode:     c["Language_Code"],
		LastName:         c["Last_Name"],
		MiddleName:       c["MI"],
		Timezone:         c["Timezone"],
		Title:            c["Title"],
		Status:           c["Receive_Email"],
		ExternalSystemID: c["supporter_KEY"],
	}
	af := []string{
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Country",
		"PostalCode",
	}

	f := false
	for _, k := range af {
		f = f || len(c[k]) > 0
	}
	if f {
		s.Address = goengage.Address{
			AddressLine1: c["Street"],
			AddressLine2: c["Street_2"],
			City:         c["City"],
			State:        c["State"],
			Country:      c["Country"],
			PostalCode:   c["Zip"],
		}
	}

	am := map[string]string{
		"Email":      "EMAIL",
		"Phone":      "HOME_PHONE",
		"Cell_Phone": "CELL_PHONE",
		"WorkPhone":  "WORK_PHONE",
	}
	as := map[string]string{
		"Email":      "OPT_IN",
		"Phone":      "",
		"Cell_Phone": "",
		"WorkPhone":  "",
	}

	var contacts []goengage.Contact
	for _, k := range af {
		if len(c[k]) > 0 {
			contact := goengage.Contact{
				Type:   am[k],
				Value:  c[k],
				Status: as[k],
			}
			contacts = append(contacts, contact)
		}
	}
	if len(contacts) > 0 {
		s.Contacts = contacts
	}
	return &s
}
