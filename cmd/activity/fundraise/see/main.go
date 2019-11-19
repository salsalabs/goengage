package main

//Application to find and detail petition signatures.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	activity "github.com/salsalabs/goengage/pkg/activity"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeFundraiseResponse(resp activity.FundraiseResponse) {
	fmt.Println("\nHeader")
	fmt.Printf("\tProcessingTime: %v\n", resp.Header.ProcessingTime)
	fmt.Printf("\tServerID: %v\n", resp.Header.ServerID)

	fmt.Println("\nPayload")
	fmt.Printf("\tTotal: %d\n", resp.Payload.Total)
	fmt.Printf("\tOffset: %d\n", resp.Payload.Offset)
	fmt.Printf("\tCount: %d\n", resp.Payload.Count)
	fmt.Printf("\tLength: %d\n", len(resp.Payload.Activities))

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
		fmt.Printf("\tTotalReceivedAmount: %5.2f\n", a.TotalReceivedAmount)
		fmt.Printf("\tDonationType: %v\n", a.DonationType)
		fmt.Printf("\tOneTimeAmount: %5.2f\n", a.OneTimeAmount)
		fmt.Printf("\tRecurringAmount: %5.2f\n", a.RecurringAmount)
		fmt.Printf("\tRecurringInterval: %v\n", a.RecurringInterval)
		fmt.Printf("\tRecurringCount: %d\n", a.RecurringCount)
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
	count := e.Metrics.MaxBatchSize
	offset := int32(0)
	for count > 0 {
		payload := activity.ActivityRequestPayload{
			Type:         activity.FundraiseType,
			Offset:       offset,
			Count:        count,
			ModifiedFrom: "2000-01-01T00:00:00.000Z",
		}
		rqt := activity.ActivityRequest{
			Header:  goengage.RequestHeader{},
			Payload: payload,
		}

		var resp activity.FundraiseResponse
		n := goengage.NetOp{
			Host:     e.Host,
			Method:   goengage.SearchMethod,
			Endpoint: activity.Search,
			Token:    e.Token,
			Request:  &rqt,
			Response: &resp,
		}
		err = n.Do()
		if err != nil {
			panic(err)
		}
		seeFundraiseResponse(resp)
		fmt.Printf("Payload total %5d, offset %5d, count %2d\n", resp.Payload.Total, resp.Payload.Offset, resp.Payload.Count)
		count = resp.Payload.Count
		offset = offset + count
	}
}
