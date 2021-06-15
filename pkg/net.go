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

const (
	//NapDuration is the time that we sleep to avoid 429 errors.  Testing shows
	//that 10 seconds is a good minimum.  Napping for two seconds repeats at least
	//five times.  We're not going anywhere anyway -- might as well wait for a
	//logoner time.  Pleae tweak as needed.
	NapDuration = "10s"

	//Multiplier is used to decide whether or not to take a nap to avoid 429 errors.
	Multiplier = 2

	//FirstWaitDuration is the duration that we nap after the first instance of a
	//HTTP 503 error.
	FirstWaitDuration = "2s"

	//MaxWaitIterations is the number of times that we'll timme out before giving up
	//because of HTTP 503's.  Note that the sleep interval doubles every time we wait.
	//MaxWaitIterations is 2 + 4 + 8 + 16 + 32 = 64 seconds.
	//so a smaller number here is better.
	MaxWaitIterations = 5
)

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
	//Loop to handle network timeouts (HTTP 504).  A 504 error is
	//insidiousconsidering that this app originally ran inside
	//Salsa's network. No WiFi, no cable companies, no kids pulling
	//wires out of th wall.
	waitDuration, _ := time.ParseDuration(FirstWaitDuration)
	ok := false

	for i := 1; !ok && i <= MaxWaitIterations; i++ {
		resp, err := n.internal()
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusGatewayTimeout {
			m := fmt.Sprintf("Error: HTTP error %v on %v. Sleeping %v seconds, pass %d of %d.",
				resp.StatusCode, n.Endpoint, waitDuration.Seconds(), i, MaxWaitIterations)
			if n.Logger != nil {
				n.Logger.Printf("%v\n", m)
			}
			log.Println(m)
			time.Sleep(waitDuration)
			waitDuration = waitDuration + waitDuration
		} else {
			ok = true
		}
	}
	return err
}

//internal processes the request provided by NetOps.  This is here so that
//we can handle both requests and metrics in the same module.
func (n *NetOp) internal() (resp *http.Response, err error) {
	u, _ := url.Parse(n.Endpoint)
	u.Scheme = "https"
	u.Host = n.Host
	var req *http.Request

	if n.Request == nil {
		req, err = http.NewRequest(n.Method, u.String(), nil)
		if err != nil {
			return nil, err
		}
	} else {
		b, err := json.Marshal(n.Request)
		if err != nil {
			return nil, err
		}
		if n.Logger != nil {
			n.Logger.LogJSON(b)
		}
		r := bytes.NewReader(b)
		req, err = http.NewRequest(n.Method, u.String(), r)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("authToken", n.Token)
	req.Header.Set("Content-Type", ContentType)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 200 {
		m := fmt.Sprintf("engage error %v on %v: %v", resp.Status, n.Endpoint, string(b))
		return resp, errors.New(m)
	}
	if n.Logger != nil {
		n.Logger.Printf("Net: endpoint %s\nNet: response\n", n.Endpoint)
		n.Logger.LogJSON(b)
	}

	err = json.Unmarshal(b, &n.Response)
	if err != nil {
		return resp, err
	}
	return resp, nil
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
	_, err = n2.internal()
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
	return b
}
