package goengage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Search submits a request and populates a response. Note
//that Engage uses HTTP status codes to denote some error
//failures.  Search passes those back to the caller as standard
//errors containing the HTTP tatus code (e.g. "200 OK").
//
//The HTTP response is unmarshalled into n.Response.
func (n *NetOp) Search() error {

	u, _ := url.Parse(n.Fragment)
	u.Scheme = "https"
	u.Host = n.Host
	//fmt.Printf("Search:  URL is %v\n", u)

	client := &http.Client{}
	rqt := RequestBase{
		//Header:  Header{},
		Payload: n.Request,
	}
	b, err := json.Marshal(rqt)
	if err != nil {
		return err
	}
	//fmt.Printf("Search: request is %v\n", string(b))
	r := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPost, u.String(), r)
	if err != nil {
		return err
	}
	req.Header.Set("authToken", n.Token)
	req.Header.Set("Content-Type", ContentType)

	resp, err := client.Do(req)
	// resp.Header.Set("Content-Type", ContentType)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v", resp.Status)
		return errors.New(m)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, n.Response)
	return err
}

//Upsert cause Engage to add/update a supporter.  If the supporter's
// ID and Email are not in the database, then Engage inserts a new
//supporter.  If either are in the database, then Engage updates the
//supporter.
//
//Note that Engage uses HTTP status codes to denote some error
//failures.  Search passes those back to the caller as standard
//errors containing the HTTP tatus code (e.g. "200 OK").
//
//The HTTP response is unmarshalled into n.Response.
func (n *NetOp) Upsert() error {

	u, _ := url.Parse(n.Fragment)
	u.Scheme = "https"
	u.Host = n.Host
	//fmt.Printf("Search:  URL is %v\n", u)

	client := &http.Client{}
	rqt := RequestBase{
		//Header:  Header{},
		Payload: n.Request,
	}
	b, err := json.Marshal(rqt)
	if err != nil {
		return err
	}
	fmt.Printf("\nSearch: request is %v\n\n", string(b))
	r := bytes.NewReader(b)

	req, err := http.NewRequest(http.MethodPut, u.String(), r)
	if err != nil {
		return err
	}
	req.Header.Set("authToken", n.Token)
	req.Header.Set("Content-Type", ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v", resp.Status)
		return errors.New(m)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(b, n.Response)
	return err
}

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

	f := false
	af := []string{
		"AddressLine1",
		"AddressLine2",
		"City",
		"State",
		"Country",
		"PostalCode",
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
