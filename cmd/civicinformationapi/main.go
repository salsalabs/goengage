package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Request contains the parameters that are sent to the API.
// See https://developers.google.com/civic-information/docs/v2/representatives/representativeInfoByAddress
type Request struct {
	Address        string
	IncludeOffices bool
	Levels         []string
	Roles          []string
}

// Response is returned by Google's
// Civic Information API when sending a requiest for representative
// information by address.
// See https://developers.google.com/civic-information/docs/v2/representatives/representativeInfoByAddress
type Response struct {
	Kind            string `json:"kind"`
	NormalizedInput struct {
		LocationName string `json:"locationName"`
		Line1        string `json:"line1"`
		Line2        string `json:"line2"`
		Line3        string `json:"line3"`
		City         string `json:"city"`
		State        string `json:"state"`
		Zip          string `json:"zip"`
	} `json:"normalizedInput"`
	Divisions struct {
		Key struct {
			Name          string   `json:"name"`
			AlsoKnownAs   []string `json:"alsoKnownAs"`
			OfficeIndices []string `json:"officeIndices"`
		} `json:"key"`
	} `json:"divisions"`
	Offices []struct {
		Name       string   `json:"name"`
		DivisionID string   `json:"divisionId"`
		Levels     []string `json:"levels"`
		Roles      []string `json:"roles"`
		Sources    []struct {
			Name     string `json:"name"`
			Official string `json:"official"`
		} `json:"sources"`
		OfficialIndices []string `json:"officialIndices"`
	} `json:"offices"`
	Officials []struct {
		Name    string `json:"name"`
		Address []struct {
			LocationName string `json:"locationName"`
			Line1        string `json:"line1"`
			Line2        string `json:"line2"`
			Line3        string `json:"line3"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
		} `json:"address"`
		Party    string   `json:"party"`
		Phones   []string `json:"phones"`
		Urls     []string `json:"urls"`
		PhotoURL string   `json:"photoUrl"`
		Emails   []string `json:"emails"`
		Channels []struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		} `json:"channels"`
	} `json:"officials"`
}

//Do submits the request and returns a slice of responses.
func Do(req Request) (resp *Response, err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return nil, nil
}

func main() {

}
