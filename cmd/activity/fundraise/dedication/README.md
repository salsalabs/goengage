# Dedication report

Go application to retrieve donations that have a dedication (in memory of or in honor of).
The user provides an Engage token and a date range.  The app finds the donations with dedications and writes selected information to a CSV.

The default date range is Monday to Sunday of last week.  See the Usage section for more information about using other date ranges.

## Prerequisites

1. A current version of Go.  There are lots of articles on the web about
installing Go.  The official installation steps can be found by [clicking here](https://golang.org/doc/install).

1. An [Engage API token](https://help.salsalabs.com/hc/en-us/articles/224470007-Salsa-Engage-Integration-API-Overview).

Note that Go *requires* a directory structure in your home directory.

```
$HOME
  + bin
  + pkg
  + src
```
  Make sure that exists before you start using Go.

## Installation

This package is part of the [GoEngage package on Github](https://github.com/salsalabs/goengage).
The package can be found in the `cmd/activity/fundraise/dedication` directory.

## Operation

1. Open a terminal window.
1. Navigate to the `goengage` directory.
```bash
cd ~/go/src/github/salsalabs.com/goengage
go run cmd/activity/fundraise/dedication/main.go [options]
 ```

### Usage

Type this to see the usage statement.
```bash
go run cmd/activity/fundraise/dedication/main.go --help
```
```
usage: dedications --login=LOGIN [<flags>]

Write dedications to a CSV

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN                  YAML file with API token
  --startDate="2020-12-14"       Start date, YYYY-MM-YY, default is Monday of last week at midnight
  --endDate="2020-12-20"         End date, YYYY-MM-YY, default is the most recent Monday at midnight
  --timezone="America/New_York"  Client's timezone, defaults to EST/EDT
```

### Command-line arguments

|Argument|Description|
|--------|-----------|
|login| LOGIN is a yaml filename containing the API token.  More on that below|
|startDate | Start of the date range for this report.  `startDate` must be formatted as "YYYY-MM-DD".  The default start date is "Monday a week ago".|
|endDate | End of the date range for this report.  `endDate` must be formatted as "YYYY-MM-DD".  The default is 7 days after `startDate`. |
|timeZone|The official timezone designation for the client.  The defualt is US Eastern.  You can more timezone names by [clickinging here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).|

### YAML file

This application needs an API token.  Tokens are are really long and typically hard to remember.  This application stores that problem by expecting a YAML file contining the API token.  Here's an example.

```yaml
token: 82874335-aedd-4f9e-918d-8ad78088a811
```
If you've stored your token in `company.yaml`, then you'll need to use a command like this to start the deduplication report.

```bash
go run cmd/activity/fundraise/dedication/main.go --login company.yaml
```

## Outputs
### Console

Here's a sample of the console output.

```bash

2020/12/23 15:48:46 WaitForReaders: Waiting for 5 readers
2020/12/23 15:48:46 ReadActivities-2: begin
2020/12/23 15:48:46 ReadActivities-1: begin
2020/12/23 15:48:46 ReadActivities-3: begin
2020/12/23 15:48:46 ReadActivities-4: begin
2020/12/23 15:48:46 WriteCSV: begin
2020/12/23 15:48:46 ReadActivities-0: begin
2020/12/23 15:48:46 ReportFundraising: processing 164 FUNDRAISE records
2020/12/23 15:48:46 ReportFundraising: waiting for terminations
2020/12/23 15:48:46 ReadActivities-4: offset    150,   1 of  14 passed
2020/12/23 15:48:46 ReadActivities-4: end
2020/12/23 15:48:46 WaitForReaders: Waiting for 4 readers
2020/12/23 15:48:46 ReadActivities-2: offset      0,   1 of  50 passed
2020/12/23 15:48:46 ReadActivities-2: end
2020/12/23 15:48:46 WaitForReaders: Waiting for 3 readers
2020/12/23 15:48:48 ReadActivities-1: offset     50,   0 of  50 passed
2020/12/23 15:48:48 ReadActivities-1: end
2020/12/23 15:48:48 WaitForReaders: Waiting for 2 readers
2020/12/23 15:48:48 ReadActivities-3: offset    100,   0 of  50 passed
2020/12/23 15:48:48 ReadActivities-3: end
2020/12/23 15:48:48 WaitForReaders: Waiting for 1 readers
2020/12/23 15:48:49 ReadActivities-0: end
2020/12/23 15:48:49 WaitForReaders: done
2020/12/23 15:48:49 WriteCSV: done
2020/12/23 15:48:49 ReportFundraising done
```
### CSV output

Output goes into `dedications.csv` in the current directory.  Here's a sample.

```
PersonName,PersonEmail,AddressLine1,AddressLine2,City,State,Zip,TransactionDate,Amount,DedicationType,Dedication
John Cheeseburger,john@cheeseburger.com.com,,,,,,2020-12-19 21:08:03.533 +0000 UTC,51.69,IN_HONOR_OF,The Cheeseburger Family
Anne Souvlaki,anne@Souvlaki.com,,,,,,2020-12-15 17:41:28.178 +0000 UTC,1912.50,IN_HONOR_OF,Our wonderful Souvlaki family
```

## Advanced usage

Typically, going to the `goengage` directory to run this app can be time-consuming.  A good way to get around that is to create an executable file.  If that appeals to you, then use this command.

```bash
go build -o ~/go/bin/fundraising_dedication cmd/activity/fundraise/dedication/main.go
```

The output will be an executable in `go/bin` in your home directory.  Add that directory to the PATH list that your OS uses and you'll be able to invoke the program with a command like this.

```bash
fundraise_dedication --login company.yaml
```

## Questions?  Comments?

Use the [GitHub issues page](https://github.com/salsalabs/goengage/issues) to report problems, ask questions or make comments. Please don't bother the nice folks at Salsalabs.  This is their nesting season and they will bite intruders.
