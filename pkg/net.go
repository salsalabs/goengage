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
	Endpoint string
	Token    string
	Request  interface{}
	Response interface{}
	Logger   *UtilLogger
}

//Do is a generic API request/response handler.  Uses the contents of
//the provided NetOp to send a request.  Parses the response back into
//the NetOp's reply.  The response in NetOp describes the complete returned
//package (fields, header, payload).
//
//Note that Engage uses HTTP status codes to denote some error
//failures.  Do passes those back to the caller as standard
//errors containing the HTTP status code (e.g. "200 OK") and the
//response body, which usually contains enlightenment about the
//error.
func (n *NetOp) Do() error {
	//Prep a request if it is provided.  Typically it is, but may not
	//be needed for some Engage API calls.  Newbie note: r is automatically
	//nil.

	u, _ := url.Parse(n.Endpoint)
	u.Scheme = "https"
	u.Host = n.Host
	var req *http.Request
	var err error

	if n.Request == nil {
		req, err = http.NewRequest(n.Method, u.String(), nil)
		if err != nil {
			return err
		}
	} else {
		b, err := json.Marshal(n.Request)
		if err != nil {
			return err
		}
		if n.Logger != nil {
			n.Logger.LogJSON(b)
		}
		r := bytes.NewReader(b)
		req, err = http.NewRequest(n.Method, u.String(), r)
		if err != nil {
			return err
		}
	}
	req.Header.Set("authToken", n.Token)
	req.Header.Set("Content-Type", ContentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v: %v", resp.Status, string(b))
		return errors.New(m)
	}
	if n.Logger != nil {
		n.Logger.LogJSON(b)
	}
	err = json.Unmarshal(b, &n.Response)
	if err != nil {
		return err
	}
	return nil
}
