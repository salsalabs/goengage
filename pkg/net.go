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

//NetOp is the wrapper for calls to Engage.  Here to keep
//call complexity down.
type NetOp struct {
	Host     string
	Method   string
	Fragment string
	Token    string
	Request  interface{}
	Response interface{}
}

//APIResponse is returned by Engage.  Payload is generally JSON and needs
//to be parsed into a struct to be useful.
type response struct {
	ID        string
	Timestamp string
	Header    struct {
		ProcessingTime int32  `json:"processingTime"`
		ServerID       string `jsin:"serverId"`
	}
	Payload interface{}
}

//Do is the Generic API request/response handler.  Fills in the Response
//for the provided NetOp object.
//
//Note that Engage uses HTTP status codes to denote some error
//failures.  Do passes those back to the caller as standard
//errors containing the HTTP status code (e.g. "200 OK") and the
//response body, which usually contains enlightenment about the
//error.
func (n *NetOp) Do() error {
	rqt := RequestBase{
		//Header:  Header{},
		Payload: n.Request,
	}
	b, err := json.Marshal(rqt)
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)

	u, _ := url.Parse(n.Fragment)
	u.Scheme = URLScheme
	u.Host = n.Host

	client := &http.Client{}
	req, err := http.NewRequest(n.Method, u.String(), r)
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
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("%v: %v", resp.Status, string(b))
		return errors.New(m)
	}
	a := response{
		Payload: &n.Response,
	}
	err = json.Unmarshal(b, &a)
	return err
}
