package goengage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//NapDuration is the time that we sleep to avoid 429 errors.
const NapDuration = "2s"

//Multiplier is used to decide whether or not to take a nap to avoid 429 errors.
const Multiplier = 2

//NetOp is the wrapper for calls to Engage.  Here to keep
//call complexity down.
type NetOp struct {
	Host     string
	Token    string
	Method   string
	Endpoint string
	Request  interface{}
	Response interface{}
	Logger   *UtilLogger
	Metrics  *Metrics
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
func (n *NetOp) Do() (err error) {
	//Avoid 429 errors by napping to build up available record slots.
	if n.Metrics == nil {
		err = n.Currently()
		if err != nil {
			return err
		}
	}
	d, _ := time.ParseDuration(NapDuration)
	for !n.Enough() {
		log.Printf("NetOp.Do: napping %v\n", d)
		time.Sleep(d)
		n.Currently()
	}
	err = n.internal()
	return err
}

//internal processes the request provided by NetOps.  This is here so that
//we can handle both requests and metrics in the same module.
func (n *NetOp) internal() (err error) {
	u, _ := url.Parse(n.Endpoint)
	u.Scheme = "https"
	u.Host = n.Host
	var req *http.Request

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
		fmt.Printf("Metrics are %+v\n", n.Metrics)
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Metrics are %+v\n", n.Metrics)
		return err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("Metrics are %+v\n", n.Metrics)
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

//Currently returns the current metrics without modifying the NetOp object.
func (n *NetOp) Currently() (err error) {
	var resp MetricsResponse
	n2 := NetOp{
		Host:     n.Host,
		Endpoint: MetricsCommand,
		Method:   http.MethodGet,
		Token:    n.Token,
		Request:  nil,
		Response: &resp,
	}
	err = n2.internal()
	if err != nil {
		return err
	}
	n.Metrics = &resp.Payload
	return nil
}

//Enough returns true if there are enough CurrentBatchSize slots to cover
//MaxBatchSize.
func (n *NetOp) Enough() bool {
	b := n.Metrics.CurrentRateLimit > Multiplier*n.Metrics.MaxBatchSize
	log.Printf("NetOp.Do: Have: %3d, Want: %3d, Enough? %v, %5v\n",
		n.Metrics.CurrentRateLimit,
		n.Metrics.MaxBatchSize,
		b,
		n.Endpoint)
	return b
}
