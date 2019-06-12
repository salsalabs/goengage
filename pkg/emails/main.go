package emails

const (
	//CommSeries allows the request to retrieve email series.
	CommSeries = "CommSeries"
	//Email allows the requst to retrieve email blasts.
	Email = "Email"
)

//IndividualReq requests email statistics for a supporter. Note that
//Cursor is used for paging and it's controlled by this library.
type IndividualReq struct {
	Cursor    string
	ID        string
	Type      string
	ContentID string
}

//EmailReq requests statistics for email blasts in a date range.
//Note that dates are ISO_8601 formatted String with a GMT timezone.
//Sample: "2020-01-05T12:34:56.000Z"
type EmailReq struct {
	PublishedFrom string
	PublishedTo   string
	Type          string
	Offset        int32
	Count         int32
}
