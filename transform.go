package goengage

import "strconv"

//SupXform transforms a map of strings into a supporter record.
func SupXform(c map[string]string) Supporter {
	s := Supporter{
		FirstName:        c["First_Name"],
		LanguageCode:     c["Language_Code"],
		LastName:         c["Last_Name"],
		MiddleName:       c["MI"],
		Timezone:         c["Timezone"],
		Title:            c["Title"],
		Status:           c["Receive_Email"],
		ExternalSystemID: c["supporter_KEY"],
	}
	if len(c["Receive_Email"]) > 0 {
		i, _ := strconv.ParseInt(c["Receive_Email"], 0, 64)
		if i > 0 {
			s.Status = "Subscribed"
		} else {
			s.Status = "Unsubscribed"
		}
	}

	f := false
	af := []string{
		"Street",
		"Street_2",
		"City",
		"State",
		"Country",
		"Zip",
	}
	for _, k := range af {
		f = f || len(c[k]) > 0
	}
	if f {
		s.Address = Address{
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
		"Work_Phone": "WORK_PHONE",
	}
	as := map[string]string{
		"Email":      "OPT_IN",
		"Phone":      "",
		"Cell_Phone": "",
		"Work_Phone": "",
	}

	var contacts []Contact
	for k, v := range am {
		if len(c[k]) > 0 {
			contact := Contact{
				Type:   v,
				Value:  c[k],
				Status: as[k],
			}
			contacts = append(contacts, contact)
		}
	}
	if len(contacts) > 0 {
		s.Contacts = contacts
	}
	return s
}
