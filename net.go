package goengage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
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

//FirstEmail returns the first email address for the provided supporter.
//Returns nil if the supporter does not have an email.  (As if...)
func FirstEmail(s Supporter) *string {
	c := s.Contacts
	if c == nil || len(c) == 0 {
		return nil
	}
	for _, x := range c {
		if x.Type == "EMAIL" {
			email := x.Value
			return &email
		}
	}
	return nil
}

//Credentials reads a YAML file with a token in it and returns the token.
func Credentials(fn string) (*EngEnv, error) {
	var c struct {
		Token string `json:"token"`
		Host  string `json:"host"`
	}
	raw, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(raw, &c)
	if err != nil {
		return nil, err
	}
	e := EngEnv{
		Token: c.Token,
		Host:  c.Host,
	}
	return &e, nil
}
