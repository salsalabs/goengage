package goengage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	//Multiplier is used to decide whether or not to take a nap to avoid 429 errors.
	Multiplier = 2

	//FirstDuration is the duration that we nap after the first instance of a
	//HTTP 504 error.
	FirstDuration = "15s"

	//MaxWaitIterations is the number of times that we'll timme out before giving up
	//because of HTTP 504's.  Note that the sleep interval doubles every time we wait.
	//MaxWaitIterations is 2 + 4 + 8 + 16 + 32 = 64 seconds.
	//so a smaller number here is better.
	MaxWaitIterations = 5
)

//NetOp is the wrapper for calls to Engage.
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

//Do is a generic API request/response handler.  Do  the contents of
//the provided NetOp to send a request.  Parses the response back
//into the NetOp's Reply.  The Response in NetOp describes the complete
//returnedpackage (fields, header, payload).
//
//Do also attempts to mitigate the effects of  HTTP 429 (too many
//requests) and 504 (network timeout) errors. Do repeats the original
//request, looking for the error condition to appear.  Each pass
//through the loop takes more time (nominally double the length of
// the last delay) an ever-increasing nap on each loop. If Do gets to
// the end of the maximum number of passes through the loop without
//relief, then Do return an error containing the HTTP status that
//caused the condition.
func (n *NetOp) Do() (err error) {
	d, _ := time.ParseDuration(FirstDuration)
	ok := false
	s := http.StatusOK

	for i := 1; !ok && i <= MaxWaitIterations; i++ {
		resp, err := n.internal()
		if err != nil {
			return err
		}
		s = resp.StatusCode
		switch s {
		case http.StatusOK:
			ok = true
		case http.StatusTooManyRequests:
			d = Delay(n, s, i, d)
		case http.StatusGatewayTimeout:
			d = Delay(n, s, i, d)
		default:
			err = fmt.Errorf("HTTP %v, %v", s, n.Endpoint)
			return err
		}
	}
	if !ok {
		err = fmt.Errorf("HTTP %v, %v", s, n.Endpoint)
		return err
	}
	return nil
}

//BotchedError returns true if the contents of the provided error message
//contains an embedded network timeout or calls-per-minute error. Side-
//effects include logging the embedded error and changing the response
//status code so that the embedded error is treated like an HTTP 429
//("too many requests").
func BotchedError(n *NetOp, resp *http.Response, s string) bool {
	ok := strings.Contains(s, "504 Gateway Time-out") ||
		strings.Contains(s, "Your per minute call rate")
	if ok {
		message := fmt.Sprintf("BotchedError: captured embedded error '%v'\n", s)
		log.Print(message)
		n.Println(message)
		resp.StatusCode = http.StatusTooManyRequests
	}
	return ok
}

//Delay displays the current HTTP status, takes a nap, and returns
//the next nap interval.
func Delay(n *NetOp, statusCode int, pass int, duration time.Duration) time.Duration {
	m := fmt.Sprintf("Delay: HTTP error %v on %v. Sleeping %v seconds, pass %d of %d.",
		statusCode, n.Endpoint, duration.Seconds(), pass, MaxWaitIterations)
	log.Println(m)
	n.Println(m)
	time.Sleep(duration)
	duration = duration * Multiplier
	return duration
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
	if n.Logger != nil {
		n.Println(fmt.Sprintf("Do: Endpoint is %v", n.Endpoint))
		n.Println(fmt.Sprintf("Do: method is %v", n.Method))
	}

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
	n.LogJSON(b)

	//This catches embedded errors in a response body.
	s := string(b)
	if BotchedError(n, resp, s) {
		return resp, nil
	}

	err = json.Unmarshal(b, &n.Response)
	if err != nil {
		//This catches embedded errors in an error message.
		s := fmt.Sprintf("%v", err)
		if BotchedError(n, resp, s) {
			return resp, nil
		}
	}
	return resp, err
}

//LogJSON writes JSON to the Logger for the provided byte slice.
func (n *NetOp) LogJSON(b []byte) {
	if n.Logger != nil {
		n.Logger.LogJSON(b)
	}
}

//Println writes to the provided Logger if it exists.
//Message should not have a trailing newline...
func (n *NetOp) Println(m string) {
	if n.Logger != nil {
		n.Logger.Printf("%v\n", m)
	}
}
