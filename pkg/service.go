package goengage

import (
	"fmt"
	"sync"
)

//Service provides the basic tools to read and filter records then
//write them to a CSV file.
type Service interface {
	//ActivityType returns the kind of activity being read.
	ActivityType() string
	//Filter returns true if the record should be used.
	Filter() bool
	//Headers returns column headers for a CSV file.
	Headers() []string
	//Line returns a list of strings to go in to the CSV file.
	Line() []string
	//Readers returns the number of readers to start.
	Readers() int
	//Filename returns the CSV filename.
	Filename() string
}

//MaxRecords returns the maximum number of activity records
//of a particular type.
func MaxRecords(e *Environment, s Service) (int32, error) {
	payload := ActivityRequestPayload{
		Type:         s.ActivityType(),
		Offset:       0,
		Count:        0,
		ModifiedFrom: "2000-01-01T00:00:00.000Z",
	}
	rqt := ActivityRequest{
		Header:  RequestHeader{},
		Payload: payload,
	}
	var resp BaseResponse
	n := NetOp{
		Host:     e.Host,
		Method:   SearchMethod,
		Endpoint: SearchActivity,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	err := n.Do()
	if err != nil {
		return int32(0), err
	}
	return int32(resp.Payload.Total), err
}

// ReportFundraising on a Service by reading all records, filtering, then
// writing survivors to a CSV file.
func ReportFundraising(e *Environment, s Service) (err error) {

	serviceChan := make(chan Service, 100)
	doneChan := make(chan bool)
	offsetChan := make(chan int32, 100)
	var wg sync.WaitGroup

	//Start the reader waiter.
	go (func(e *Environment, s Service, c chan Service, done chan bool, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WaitForReaders(e, s, c, done)
	})(e, s, serviceChan, doneChan, &wg)

	//Start the CSV writer.
	go (func(e *Environment, s Service, c chan Service, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		WriteCSV(e, s, c)
	})(e, s, serviceChan, &wg)

	//Start the readers.
	for i := 0; i < s.Readers(); i++ {
		go (func(e *Environment, s Service, i int, offset chan int32, c chan Service, done chan bool, wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			ReadActivities(e, s, i, offsetChan, serviceChan, doneChan)
		})(e, s, i, offsetChan, serviceChan, doneChan, &wg)
	}

	maxRecords, err := MaxRecords(e, s)
	fmt.Printf("%d %s records\n", maxRecords, FundraiseType)
	return err
}

//ReadActivities retrieves activity records from Engage, filters them,
//then writes them to the Service channel. The offset channel tells
//us where to start reading.  When no items are available from the
//offset channel, we'll write a true to the done channel.
func ReadActivities(e *Environment, s Service, i int, offsetChan chan int32, serviceChan chan Service, doneChan chanBool) {

}

//WriteCSV waits for Service records to appear on the queue, then
//Writes them to a CSV file.
func WriteCSV(e *Environment, s Service, c chan Service) {

}

//WaitForReaders waits for readers to send to a done channel.
//The number of readers is specified in the provided Service.
//Closes the csv channel when all readers are done.
func WaitForReaders(e *Environment, s Service, c chan Service, done chan bool) {
	count := s.Readers()
	for count > 0 {
		_, ok := <-c
		if !ok {
			break
		}
		count--
	}
	close(c)
}
