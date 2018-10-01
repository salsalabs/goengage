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
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v: %v", resp.Status, string(b))
		return errors.New(m)
	}
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
	fmt.Printf("Upsert: request is %v\n\n", string(b))
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
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v: %v", resp.Status, string(b))
		return errors.New(m)
	}
	err = json.Unmarshal(b, n.Response)
	return err
}
