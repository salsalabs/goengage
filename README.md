# goengage
Engage API using Go.

## Summary
This is the first pass at creating a Go version of the Engage API.  The basic data
structures all appear in the root directory as do network functions and utilities.
The "cmd" directory contains applications that use the basic API in the root directory.

## Installation

```bash
go get github.com/salsalabs/goengage
go get ./...
```

## Typical Usage
```go
import (
    "github.com/salsalabs/goengage"
    kingpin "gopkg.in/alecthomas/kingpin.v2"
)
func main() {
	var (
		app   = kingpin.New("activity-search", "A command-line app to search for supporters added by activities.")
		login = app.Flag("login", "YAML file with API token").Required().String()
	)
	app.Parse(os.Args[1:])
	e, err := goengage.Credentials(*login)
	if err != nil {
		panic(err)
    }
    	m, err := e.Metrics()
	if err != nil {
		panic(err)
	}
	rqt := goengage.SegSearchRequest{
		Offset:       0,
		Count:        m.MaxBatchSize,
		MemberCounts: !*fast,
	}
	var resp goengage.SegSearchResult
	n := goengage.NetOp{
		Host:     e.Host,
		Fragment: goengage.SegSearch,
		Token:    e.Token,
		Request:  &rqt,
		Response: &resp,
    }
    
    err = n.Search()
    if err != nil {
        panic(err)
    }
    //Internate through the items in `resp` to do stuff.
}
```

## Access

Unlike Salsa Classic, Engage does not have a login step.  Each API call must be sent to the correct host accompanies by an API Token.

The easiest way to do that is to put the host and the token into a YAML file.   The applications
provided with the API use YAML files to provide that information.

Here's a sample for a production Engage instance.

```yaml
Host: api.salsalabs.com
Token: mary-had-little-lamb-its-fleece-was-white-as-snow
```
Here's a sample for an instance of Engage that's on Salsa's internal UAT site.

```yaml
Host: hq.uat.igniteaction.org
Token: nowisthetimefor_a_quickbrownfox_to_jumpoveralazydog
```
Please read [the Engage documentation](https://help.salsalabs.com/hc/en-us/sections/205407008-API-Engage-Integration) to learn
more about API hosts and tokens.

## Applications included
This is a partial list of the applications that are distributed with the Engage API.

### `cmd/activity/added_supporters`

This application was written for a client that used the Classic `Source_Details` field.  
`Source_Details` holds the URL of the page that created a supporter.  Engage does not have
an equivalent.

This application was written to find activities where a supporter was created.  There's
not yet a nice little flag to says "This action created a supporter".  As an alternative, 
the app compares the time that the activity was created and the time that the supporter was
created.  If they are less than a second apart, then the app presumes that the action created
the supporter.

Supporters that were created by actions are written to disk in CSV format.  The filename
is hardcoded as `supporter_page.csv`.  

Errors are noisy and fatal.

Here's a sample of the help for the application.

```bash

go run cmd/activity/added_suporters/main.go --help

usage: activity-search --login=LOGIN [<flags>]

A command-line app to search for supporters added by activities.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
```

### `cmd/activity/see`

An application that reads through all activities and displays basic
activity information.  The app is inteded to be a starting place for
a developer that wants to learn how to peruse Engage activities.

Usage:

```bash
go run cmd/activity/see/main.go --help
usage: activity-search --login=LOGIN [<flags>]

A command-line app to search for supporters added by activities.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
  ```

  Sample output:
  ```
  018/09/27 17:12:29 Main: napping and then waiting.
2018/09/27 17:12:29 Lookup: start
2018/09/27 17:12:29 Merge: start
2018/09/27 17:12:30 Drive: max size is 20, we're using 20
2018/09/27 17:12:31 Drive: read 20 activities from offset 0
2018/09/27 17:12:31 Drive: read 20 activities from offset 20
layla.maher5235@uatauto.ignite.net       2015-12-21     TARGETED_LETTER Lets save the trees     -2.334s
rio.rose3164@uatauto.ignite.net          2015-12-21     TARGETED_LETTER Lets save the trees     49m36.319s
zita.morrison7275@uatauto.ignite.net     2015-12-21     TARGETED_LETTER Lets save the trees     46m42.99s
robin.irwin2312@uatauto.ignite.net       2015-12-21     TARGETED_LETTER Lets save the trees     51m29.84s
aniya.rubio4883@uatauto.ignite.net       2015-12-21     TARGETED_LETTER Lets save the trees     52m47.861s
shannon.prince3455@uatauto.ignite.net    2015-12-21     TARGETED_LETTER Lets save the trees     47m43.596s
ana.carroll1141@uatauto.ignite.net       2015-12-22     SUBSCRIBE       Follow us       -1.729s
```

### `cmd/metrics/main`

This app retrieves the current metrics from Engage.  You can learn 
more about metrics by [clicking here](https://help.salsalabs.com/hc/en-us/articles/224470007-Getting-Started).

Usage:
```bash
go run cmd/metrics/main.go --help
usage: metrics --login=LOGIN [<flags>]

A command-line app to display the current Engage metrics for a token.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
  ```

  Sample Output:
  ```
  Setting                        Value
------------------------------ -------------------------
RateLimit                      300
MaxBatchSize                   20
CurrentRateLimit               300
TotalAPICalls                  2835
LastAPICall                    2018-09-27T22:12:36.704Z
TotalAPICallFailures           0
LastAPICallFailure
SupporterRead                  6542
SupporterAdd                   1772
SupporterUpdate                69
SupporterDelete                0
ActivityEvent                  0
ActivitySubscribe              4664
ActivityFundraise              86
ActivityTargetedLetter         925
ActivityPetition               1094
ActivitySubscriptionManagement 37
```

### `cmd/segment/search`

An application that scans the database for segments (groups).  Each
line contains a selection of information that's avaiable.

Usage:
```bash
go run cmd/segment/search/main.go --help
usage: see-segments --login=LOGIN [<flags>]

A command-line app to search for segments.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
  --fast         Don't show number of members
  ```

Here's a sample of the default (slow) output.

```
SegmentID                            Name                                     Type       Members ExtID    Description
------------------------------------ ---------------------------------------- ---------- ------- -------- -------------------------
Reading 20 from 0
4ab2f622-1e14-4c70-ae7e-9650de992f42 Donors over $50                          CUSTOM         877
79ebfbd8-0382-4f5f-80ad-971a85de6b06 Has Never Made a Donation                DEFAULT     105301
306a5324-8788-43b3-85b5-fe01886bc02e Donor Subscribers                        DEFAULT        913
2b9549d6-3848-4528-a470-f0ab98687f36 Emerging Donors                          DEFAULT        258
a6354b29-43bb-4bbe-85a8-09d10248f9c3 Source is a Petition Form                DEFAULT       1881
4cbafb62-630f-4dcd-816e-10ef5b2fa018 Social Subscribers                       DEFAULT         24
f115c126-0577-49ce-82c2-9036356445f5 Dog People                               CUSTOM           0          People who indicate that they own dogs.  They may also own cats, but they do own a dog.

```

Here's a sample of the "fast" output.

```
SegmentID                            Name                                     Type       Members ExtID    Description
------------------------------------ ---------------------------------------- ---------- ------- -------- -------------------------
Reading 20 from 0
4ab2f622-1e14-4c70-ae7e-9650de992f42 Donors over $50                          CUSTOM           0
79ebfbd8-0382-4f5f-80ad-971a85de6b06 Has Never Made a Donation                DEFAULT          0
306a5324-8788-43b3-85b5-fe01886bc02e Donor Subscribers                        DEFAULT          0
2b9549d6-3848-4528-a470-f0ab98687f36 Emerging Donors                          DEFAULT          0
a6354b29-43bb-4bbe-85a8-09d10248f9c3 Source is a Petition Form                DEFAULT          0
4cbafb62-630f-4dcd-816e-10ef5b2fa018 Social Subscribers                       DEFAULT          0
f115c126-0577-49ce-82c2-9036356445f5 Dog People                               CUSTOM           0          People who indicate that they own dogs.  They may also own cats, but they do own a dog.
```
# `cmd/supporter/search`

Application that exercises the supporter search function in Engage API.

Usage:

```bash
go run cmd/supporter/search/main.go --help
usage: activity-search --login=LOGIN [<flags>]

A command-line app to see all supporters.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
```

Sample output:

```
go run cmd/supporter/search/main.go --login logins/sandbox.yaml

Searching from offset 0
Read 20 supporters from offset 0
Zayla                Hopkins              zulu@yankee.xray
Ria                  Kelley               able8411@ugregory.ba
DonorPro             System               baker@shotel.bb
Allen                Leonard              charlie@india.bc
test                 test                 delta@joliet.bd
Debbie               Williams             eagle@kilo.be
Salsa Staff Test     Salsa Staff Test     foxtrot@lima.bes
```

# `cmd/supporter/see`
This application accepts an email address and displays supporters
that have that address.

Usage:
```bash
go run cmd/supporter/see/main.go --help
usage: see-supporter --login=LOGIN --email=EMAIL [<flags>]

A command-line app to to show supporters for an email.

Flags:
  --help         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN  YAML file with API token
  --email=EMAIL  Email address to look up
  ```
  Output:
  ```
  go run cmd/supporter/see/main.go --login logins/sandbox.yaml --email aleonard@salsalabs.com
{
  "SupporterID": "ea8356fc-eb91-4bde-b514-877322bd6996",
  "Result": "FOUND",
  "Title": "Dr",
  "FirstName": "test",
  "MiddleName": "t",
  "LastName": "test",
  "Suffix": "",
  "DateOfBirth": "",
  "Gender": "",
  "CreatedDate": "2016-12-14T17:36:40.698Z",
  "LastModified": "2018-09-25T17:05:09.968Z",
  "ExternalSystemID": "",
  "Address": {
    "AddressLine1": "test",
    "AddressLine2": "test",
    "City": "test",
    "State": "TX",
    "PostalCode": "78701",
    "County": "",
    "Country": "US",
    "FederalDistrict": "",
    "StateHouseDistrict": "",
    "StateSenateDistrict": "",
    "CountyDistrict": "",
    "MunicipalityDistrict": "",
    "Lattitude": 0,
    "Longitude": 0,
    "Status": "OptIn"
  },
  "Contacts": [
    {
      "Type": "WORK_PHONE",
      "Value": "512.555.1212",
      "Status": "",
      "Errors": null
    },
    {
      "Type": "HOME_PHONE",
      "Value": "512-555-1313",
      "Status": "",
      "Errors": null
    },
    {
      "Type": "EMAIL",
      "Value": "aleonard@salsalabs.com",
      "Status": "OPT_IN",
      "Errors": null
    }
  ],
  "CustomFieldValues": []
}
```

## License

See the LICENSE file in this directory. 

## Questions

Use the [Issues](https://github.com/salsalabs/goengage/issues) link
at the top of this page to report issues.  Please don't waste your time
by contacting Salsalabs Support.
