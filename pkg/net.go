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
}

//Do is a generic API request/response handler.  Uses the contents of
//the provided NetOp to send a request.  Parses the response back into
//the NetOp's reply.
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
	var r *bytes.Reader
	var err error
	if n.Request != nil {
		rqt := Request{
			Header:  RequestHeader{},
			Payload: n.Request,
		}
		b, err := json.Marshal(rqt)
		if err != nil {
			return err
		}
		r = bytes.NewReader(b)
	}

	u, _ := url.Parse(n.Endpoint)
	u.Scheme = "https"
	u.Host = n.Host

	client := &http.Client{}
	var req *http.Request
	if n.Request == nil {
		//'r' is a concrete instantiation.  Setting it to nil is not the
		//same as passing a nil, apparently.  Interesting, no?
		req, err = http.NewRequest(n.Method, u.String(), nil)
	} else {
		req, err = http.NewRequest(n.Method, u.String(), r)
	}
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Body: %v\n", string(b))
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v: %v", resp.Status, string(b))
		return errors.New(m)
	}
	var x = Response{
		Payload: n.Response,
	}
	err = json.Unmarshal(b, &x)
	fmt.Printf("Unmarshalled response: %+v\n", x)
	fmt.Printf("Errors: %v\n", x.Errors)
	if len(x.Errors) > 0 {
		e := x.Errors[0]
		fmt.Printf("Engage error: %+v\n", e)
		m := fmt.Sprintf("Engage error %+v\n", e)
		return errors.New(m)
	}
	return err
}
