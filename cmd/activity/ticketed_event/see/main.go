package main

//Application to find and detail petition signatures.
import (
	"fmt"
	"os"

	goengage "github.com/salsalabs/goengage/pkg"
	activity "github.com/salsalabs/goengage/pkg/activity"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func seeTicketedEventResponse(resp activity.TicketedEventResponse) {
	fmt.Println("\nHeader")
	fmt.Printf("\tProcessingTime: %v\n", resp.Header.ProcessingTime)
	fmt.Printf("\tServerID: %v\n", resp.Header.ServerID)

	fmt.Println("\nPayload")
	fmt.Printf("\tTotal: %v\n", resp.Payload.Total)
	fmt.Printf("\tOffset: %v\n", resp.Payload.Offset)
	fmt.Printf("\tCount: %v\n", resp.Payload.Count)
	fmt.Printf("\tLength: %v\n", len(resp.Payload.Activities))

	fmt.Println("\nTicketedEvents")
	for i, e := range resp.Payload.Activities {
		fmt.Printf("\n\tTicketedEvent %d\n", i)
		fmt.Printf("\tActivityID: %v\n", e.ActivityID)
		fmt.Printf("\tActivityFormName: %v\n", e.ActivityFormName)
		fmt.Printf("\tActivityFormID: %v\n", e.ActivityFormID)
		fmt.Printf("\tSupporterID: %v\n", e.SupporterID)
		fmt.Printf("\tActivityDate: %v\n", e.ActivityDate)
		fmt.Printf("\tActivityType: %v\n", e.ActivityType)
		fmt.Printf("\tLastModified: %v\n", e.LastModified)
		fmt.Printf("\tDonationID: %v\n", e.DonationID)
		fmt.Printf("\tTotalReceivedAmount: %v\n", e.TotalReceivedAmount)
		fmt.Printf("\tOneTimeAmount: %v\n", e.OneTimeAmount)
		fmt.Printf("\tDonationType: %v\n", e.DonationType)
		fmt.Printf("\tAccountType: %v\n", e.AccountType)
		fmt.Printf("\tAccountNumber: %v\n", e.AccountNumber)
		fmt.Printf("\tAccountExpiration: %v\n", e.AccountExpiration)
		fmt.Printf("\tAccountProvider: %v\n", e.AccountProvider)
		fmt.Printf("\tPaymentProcessorName: %v\n", e.PaymentProcessorName)
		fmt.Printf("\tActivityResult: %v\n", e.ActivityResult)

		fmt.Printf("\n\tTransactions")
		for j, x := range e.Transactions {
			fmt.Printf("\n\t\tTransaction %d\n", j)
			fmt.Printf("\t\tTransactionID: %v\n", x.TransactionID)
			fmt.Printf("\t\tType: %v\n", x.Type)
			fmt.Printf("\t\tReason: %v\n", x.Reason)
			fmt.Printf("\t\tDate: %v\n", x.Date)
			fmt.Printf("\t\tAmount: %v\n", x.Amount)
			fmt.Printf("\t\tDeductibleAmount: %v\n", x.DeductibleAmount)
			fmt.Printf("\t\tFeesPaid: %v\n", x.FeesPaid)
			fmt.Printf("\t\tGatewayTransactionID: %v\n", x.GatewayTransactionID)
			fmt.Printf("\t\tGatewayAuthorizationCode: %v\n", x.GatewayAuthorizationCode)
		}
		fmt.Printf("\n\tTickets")
		for j, t := range e.Tickets {
			fmt.Printf("\n\t\tTicket %d\n", j)
			fmt.Printf("\t\tTicketID: %v\n", t.TicketID)
			fmt.Printf("\t\tTicketName: %v\n", t.TicketName)
			fmt.Printf("\t\tTransactionID: %v\n", t.TransactionID)
			fmt.Printf("\t\tLastModified: %v\n", t.LastModified)
			fmt.Printf("\t\tTicketStatus: %v\n", t.TicketStatus)
			fmt.Printf("\t\tTicketCost: %v\n", t.TicketCost)
			fmt.Printf("\t\tDeductibleAmount: %v\n", t.DeductibleAmount)
			fmt.Printf("\n\t\tQuestions")
			for k, q := range t.Questions {
				fmt.Printf("\n\t\t\tQuestion %d\n", k)
				fmt.Printf("\t\t\tID: %v\n", q.ID)
				fmt.Printf("\t\t\tQuestion: %v\n", q.Question)
				fmt.Printf("\t\t\tAnswer: %v\n", q.Answer)
			}
			fmt.Printf("\n\t\tAttendees")
			for k, a := range t.Attendees {
				fmt.Printf("\n\t\t\tAttendee %d\n", k)
				fmt.Printf("\t\t\tAttendeeID: %v\n", a.AttendeeID)
				fmt.Printf("\t\t\tFirstName: %v\n", a.FirstName)
				fmt.Printf("\t\t\tType: %v\n", a.Type)
				fmt.Printf("\t\t\tStatus: %v\n", a.Status)
				fmt.Printf("\t\t\tLastName: %v\n", a.LastName)
				fmt.Printf("\t\t\tEmail: %v\n", a.Email)
				fmt.Printf("\t\t\tAdressLine1: %v\n", a.AdressLine1)
				fmt.Printf("\t\t\tAdressLine2: %v\n", a.AdressLine2)
				fmt.Printf("\t\t\tCity: %v\n", a.City)
				fmt.Printf("\t\t\tState: %v\n", a.State)
				fmt.Printf("\t\t\tPhone: %v\n", a.Phone)
				fmt.Printf("\t\t\tIsCurrentSupporter: %v\n", a.IsCurrentSupporter)
				fmt.Printf("\t\t\tLastModified: %v\n", a.LastModified)

				fmt.Printf("\n\t\t\tQuestions: %v\n", a.Questions)
				for k, q := range a.Questions {
					fmt.Printf("\n\t\t\t\tQuestion %d\n", k)
					fmt.Printf("\t\t\t\tID: %v\n", q.ID)
					fmt.Printf("\t\t\t\tQuestion: %v\n", q.Question)
					fmt.Printf("\t\t\t\tAnswer: %v\n", q.Answer)
				}

			}
		}
		fmt.Printf("\n\tPurchases")
		for j, p := range e.Purchases {
			fmt.Printf("\n\t\tPurchase %d\n", j)
			fmt.Printf("\t\tPurchaseID: %v\n", p.PurchaseID)
			fmt.Printf("\t\tTicketID: %v\n", p.TicketID)
			fmt.Printf("\t\tAttendeeID: %v\n", p.AttendeeID)
			fmt.Printf("\t\tName: %v\n", p.Name)
			fmt.Printf("\t\tCost: %v\n", p.Cost)
			fmt.Printf("\t\tQuantity: %v\n", p.Quantity)
			fmt.Printf("\t\tStatus: %v\n", p.Status)

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
	payload := activity.ActivityRequestPayload{
		Type:         activity.TicketedEventType,
		Offset:       0,
		Count:        e.Metrics.MaxBatchSize,
		ModifiedFrom: "2000-01-01T00:00:00.000Z",
	}
	rqt := activity.ActivityRequest{
		Header:  goengage.Header{},
		Payload: payload,
	}
	var resp activity.TicketedEventResponse
	n := goengage.NetOp{
		Host:     e.Host,
		Method:   goengage.SearchMethod,
		Endpoint: activity.Search,
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
	seeTicketedEventResponse(resp)
}
