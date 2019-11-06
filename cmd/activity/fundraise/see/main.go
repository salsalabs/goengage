package main

//Application to find and detail petition signatures.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeFundraiseResponse(resp goengage.FundraiseResponse) {
	fmt.Println("\nHeader")
	fmt.Printf("\tProcessingTime: %v\n", resp.Header.ProcessingTime)
	fmt.Printf("\tServerID: %v\n", resp.Header.ServerID)

	fmt.Println("\nPayload")
	fmt.Printf("\tTotal: %v\n", resp.Payload.Total)
	fmt.Printf("\tOffset: %v\n", resp.Payload.Offset)
	fmt.Printf("\tCount: %v\n", resp.Payload.Count)
	fmt.Printf("\tLength: %v\n", len(resp.Payload.Activities))

	fmt.Println("\nFundraise")
	for i, a := range resp.Payload.Activities {
		fmt.Printf("\n\tFundraise %d\n", i)
		fmt.Printf("\tActivityID: %v\n", a.ActivityID)
		fmt.Printf("\tActivityFormName: %v\n", a.ActivityFormName)
		fmt.Printf("\tActivityFormID: %v\n", a.ActivityFormID)
		fmt.Printf("\tSupporterID: %v\n", a.SupporterID)
		fmt.Printf("\tActivityDate: %v\n", a.ActivityDate)
		fmt.Printf("\tActivityType: %v\n", a.ActivityType)
		fmt.Printf("\tLastModified: %v\n", a.LastModified)
		fmt.Printf("\tDonationID: %v\n", a.DonationID)
		fmt.Printf("\tTotalReceivedAmount: %v\n", a.TotalReceivedAmount)
		fmt.Printf("\tDonationType: %v\n", a.DonationType)
		fmt.Printf("\tOneTimeAmount: %v\n", a.OneTimeAmount)
		fmt.Printf("\tRecurringAmount: %v\n", a.RecurringAmount)
		fmt.Printf("\tRecurringInterval: %v\n", a.RecurringInterval)
		fmt.Printf("\tRecurringCount: %v\n", a.RecurringCount)
		fmt.Printf("\tRecurringTransactionID: %v\n", a.RecurringTransactionID)
		fmt.Printf("\tRecurringStart: %v\n", a.RecurringStart)
		fmt.Printf("\tRecurringEnd: %v\n", a.RecurringEnd)
		fmt.Printf("\tAccountType: %v\n", a.AccountType)
		fmt.Printf("\tAccountNumber: %v\n", a.AccountNumber)
		fmt.Printf("\tAccountExpiration: %v\n", a.AccountExpiration)
		fmt.Printf("\tAccountProvider: %v\n", a.AccountProvider)
		fmt.Printf("\tPaymentProcessorName: %v\n", a.PaymentProcessorName)
		fmt.Printf("\tFundName: %v\n", a.FundName)
		fmt.Printf("\tFundGLCode: %v\n", a.FundGLCode)
		fmt.Printf("\tDesignation: %v\n", a.Designation)
		fmt.Printf("\tDedicationType: %v\n", a.DedicationType)
		fmt.Printf("\tDedication: %v\n", a.Dedication)
		fmt.Printf("\tNotify: %v\n", a.Notify)

		fmt.Printf("\tTransactions")
		for i, t := range a.Transactions {
			fmt.Printf("\n\t\tTrasaction %d\n", i)
			fmt.Printf("\t\tTransactionID: %v\n", t.TransactionID)
			fmt.Printf("\t\tType: %v\n", t.Type)
			fmt.Printf("\t\tReason: %v\n", t.Reason)
			fmt.Printf("\t\tDate: %v\n", t.Date)
			fmt.Printf("\t\tAmount: %v\n", t.Amount)
			fmt.Printf("\t\tDeductibleAmount: %v\n", t.DeductibleAmount)
			fmt.Printf("\t\tFeesPaid: %v\n", t.FeesPaid)
			fmt.Printf("\t\tGatewayTransactionID: %v\n", t.GatewayTransactionID)
			fmt.Printf("\t\tGatewayAuthorizationCode: %v\n", t.GatewayAuthorizationCode)
		}
	}
}

func main() {
	var (
		app   = kingpin.New("activity-see", "List all activities")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
	}
	rqt := goengage.ActivityRequest{
		Type:         goengage.FundraiseType,
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
		ModifiedFrom: "2010-01-01T00:00:00.000Z",
	}
	var resp goengage.FundraiseResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: goengage.ActSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
	}
	//b, _ := json.MarshalIndent(n, "", "    ")
	//fmt.Printf("NetOp: %+v\n", string(b))

	err = n.Do()
	if err != nil {
		panic(err)
	}
	//b, _ = json.MarshalIndent(rqt, "", "    ")
	//fmt.Printf("Request: %+v\n", string(b))
	//b, _ = json.MarshalIndent(resp, "", "    ")
	//fmt.Printf("Response: %+v\n", string(b))
	seeFundraiseResponse(resp)
}
