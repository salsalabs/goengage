# Dedication report

Go application to retrieve donations that have a dedication (in memory of or in honor of).

The user provides an Engage token and a date range.
The default date range is Monday to Sunday of last week.
See the Usage section for more information.

The app finds the donations with dedications and writes information to a CSV.
* Personal information (name, email, address)
* Donation information (date, amount)
* Dedication information (type, dedication, notify)
* Dedication address information (custom field in the donor's supporter record)


## Prerequisites

1. A current version of Go.  There are lots of articles on the web about
installing Go.  The official installation steps can be found by [clicking here](https://golang.org/doc/install).
1. An [Engage API token](https://help.salsalabs.com/hc/en-us/articles/224470007-Salsa-Engage-Integration-API-Overview).


## Installation

This package is part of the [GoEngage package on Github](https://github.com/salsalabs/goengage). 
Use these steps to install `goengage`.

```bash
go get github.com/salsalabs/goengage
go install github.com/salsalabs/goengage
```

The source for this package can be found in the `cmd/activity/fundraise/dedication` directory in `goengage`.

## Operation

The easiest way to run this app is to start in a console window. 

```bash
cd ~/go/src/github/salsalabs.com/goengage
go run cmd/activity/fundraise/dedication/main.go --help
```

Use `--help` shows the usage summary.

```
usage: dedications --login=LOGIN [<flags>]

Write dedications to a CSV

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --login=LOGIN                  YAML file with API token
  --startDate="2020-12-14"       Start date, YYYY-MM-YY, default is Monday of last week at midnight
  --endDate="2020-12-20"         End date, YYYY-MM-YY, default is the most recent Monday at midnight
  --timezone="America/New_York"  Client's timezone, defaults to EST/EDT
   --keys                        Export activity, donation and supporter IDs
```

### Command-line arguments

|Argument|Description|
|--------|-----------|
|login| LOGIN is a yaml filename containing the API token.  See [YAML file](#yaml-file), below.|
|startDate | Start of the date range for this report.  `startDate` must be formatted as "YYYY-MM-DD".  The default start date is "Monday a week ago".  The default date appears in the usage.|
|endDate | End of the date range for this report.  `endDate` must be formatted as "YYYY-MM-DD".  The default is 7 days after `startDate`. It, too, appears in the usage.|
|timeZone|The official timezone designation for the client.  The default is US Eastern.  You can see more timezone names by [clicking here](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).|
|keys|Use this to append the ActivityID, DonationID and SupporterID to each donation.|

### YAML file

This application needs an API token.
Tokens are are really long and typically hard to remember.
This application stores that problem by expecting a YAML file containing the API token.  Here's an example.

```yaml
token: 82874335-aedd-4f9e-918d-8ad78088a811
```
If you've stored your token in `company.yaml`, then you'll need to use a command like this to start the deduplication report.

```bash
go run cmd/activity/fundraise/dedication/main.go --login company.yaml
```

## Outputs
### Console

The application logs all status to the console.  Errors are really obvious.  Here's a sample of the console output for a successful application run.

```
2020/12/30 09:58:33 WaitForReaders: Waiting for 5 readers
2020/12/30 09:58:33 ReadActivities-1: begin
2020/12/30 09:58:33 ReadActivities-0: begin
2020/12/30 09:58:33 ReadActivities-2: begin
2020/12/30 09:58:33 ReadActivities-3: begin
2020/12/30 09:58:33 WriteCSV: begin
2020/12/30 09:58:33 ReadActivities-4: begin
2020/12/30 09:58:33 ReportFundraising: processing 366 FUNDRAISE records
2020/12/30 09:58:33 ReportFundraising: waiting for terminations
2020/12/30 09:58:37 ReadActivities-2: offset     40 of    366,   6 adds
2020/12/30 09:58:39 ReadActivities-3: offset     60 of    366,   5 adds
2020/12/30 09:58:40 ReadActivities-4: offset     80 of    366,  13 adds
2020/12/30 09:58:42 ReadActivities-4: offset    140 of    366,   9 adds
2020/12/30 09:58:44 ReadActivities-3: offset    120 of    366,  11 adds
2020/12/30 09:58:46 ReadActivities-2: offset    100 of    366,  11 adds
2020/12/30 09:58:46 ReadActivities-3: offset    180 of    366,   8 adds
2020/12/30 09:58:48 ReadActivities-2: offset    200 of    366,   7 adds
2020/12/30 09:58:51 ReadActivities-0: offset     20 of    366,  10 adds
2020/12/30 09:58:51 ReadActivities-4: offset    160 of    366,  16 adds
2020/12/30 09:58:52 ReadActivities-3: offset    220 of    366,  17 adds
2020/12/30 09:58:55 ReadActivities-4: offset    280 of    366,   8 adds
2020/12/30 09:58:55 ReadActivities-3: offset    300 of    366,  11 adds
2020/12/30 09:58:58 ReadActivities-3: offset    340 of    366,   9 adds
2020/12/30 09:58:58 ReadActivities-4: offset    320 of    366,   8 adds
2020/12/30 09:58:58 ReadActivities-3: offset    360 of    366,   2 adds
2020/12/30 09:58:58 ReadActivities-3: end
2020/12/30 09:58:58 WaitForReaders: Waiting for 4 readers
2020/12/30 09:58:59 ReadActivities-4: end
2020/12/30 09:58:59 WaitForReaders: Waiting for 3 readers
2020/12/30 09:59:01 ReadActivities-1: offset      0 of    366,   9 adds
2020/12/30 09:59:01 ReadActivities-1: end
2020/12/30 09:59:01 WaitForReaders: Waiting for 2 readers
2020/12/30 09:59:02 ReadActivities-0: offset    260 of    366,  12 adds
2020/12/30 09:59:02 ReadActivities-0: end
2020/12/30 09:59:02 WaitForReaders: Waiting for 1 readers
2020/12/30 09:59:03 ReadActivities-2: offset    240 of    366,  10 adds
2020/12/30 09:59:03 ReadActivities-2: end
2020/12/30 09:59:03 WaitForReaders: done
2020/12/30 09:59:03 WriteCSV: done
2020/12/30 09:59:03 ReportFundraising done
```
If you choose start and end dates in different months, the application will process each month separately.  Here's an example.

```
go run main.go --login ~/.logins/mules.yaml --startDate "2021-01-01" --endDate "2021-02-28"
2021/07/21 11:25:14 
2021/07/21 11:25:14 WaitForReaders: Waiting for 3 readers
2021/07/21 11:25:14 Store: begin
2021/07/21 11:25:14 ReadActivities-1: begin
2021/07/21 11:25:14 ReadActivities-0: begin
2021/07/21 11:25:14 ReadActivities-2: begin
2021/07/21 11:25:15 ReportFundraising: reporting on start time 2021-01-01T05:00:00.000Z
2021/07/21 11:25:15 ReportFundraising:              end   time 2021-02-01T04:59:59.999Z
2021/07/21 11:25:15 ReportFundraising: 23 donations
2021/07/21 11:25:15 ReportFundraising: waiting for terminations
2021/07/21 11:25:15 ReadActivities-2: end
2021/07/21 11:25:15 WaitForReaders: Waiting for 2 readers
2021/07/21 11:25:16 ReadActivities-0: offset     20 of     23,   2 adds
2021/07/21 11:25:16 ReadActivities-0: end
2021/07/21 11:25:16 WaitForReaders: Waiting for 1 readers
2021/07/21 11:25:20 ReadActivities-1: offset      0 of     23,  15 adds
2021/07/21 11:25:20 ReadActivities-1: end
2021/07/21 11:25:20 WaitForReaders: done
2021/07/21 11:25:20 Store: done
2021/07/21 11:25:20 ReportFundraising: done
2021/07/21 11:25:20 
2021/07/21 11:25:20 Store: begin
2021/07/21 11:25:20 ReadActivities-0: begin
2021/07/21 11:25:20 WaitForReaders: Waiting for 3 readers
2021/07/21 11:25:20 ReadActivities-1: begin
2021/07/21 11:25:20 ReadActivities-2: begin
2021/07/21 11:25:21 ReportFundraising: reporting on start time 2021-02-01T05:00:00.000Z
2021/07/21 11:25:21 ReportFundraising:              end   time 2021-03-01T04:59:59.999Z
2021/07/21 11:25:21 ReportFundraising: 11 donations
2021/07/21 11:25:21 ReportFundraising: waiting for terminations
2021/07/21 11:25:21 ReadActivities-2: end
2021/07/21 11:25:21 WaitForReaders: Waiting for 2 readers
2021/07/21 11:25:21 ReadActivities-1: end
2021/07/21 11:25:21 WaitForReaders: Waiting for 1 readers
2021/07/21 11:25:23 ReadActivities-0: offset      0 of     11,   9 adds
2021/07/21 11:25:23 ReadActivities-0: end
2021/07/21 11:25:23 WaitForReaders: done
2021/07/21 11:25:23 Store: done
2021/07/21 11:25:23 ReportFundraising: done
```
Note that each month appears in its own CSV.

```
ls -al *.csv
-rw-r--r--  1 aleonard  staff  2816 Jul 21 11:25 2021-01-01_dedications.csv
-rw-r--r--  1 aleonard  staff  1633 Jul 21 11:25 2021-02-01_dedications.csv

```

### CSV output

Results are stored in CSV files.  The date in the CSV filename is the first date
of the reporting period.  Here's a sample of the CSV output.

```
FirstName,LastName,PersonEmail,AddressLine1,AddressLine2,City,State,Zip,TransactionDate,DonationType,ActivityType,TransactionType,Amount,DedicationType,Dedication,Notify,DedicationAddress
Patsy,Pastry,patsy@pastry.com,273 Ramblin Man,,Grapevine,MI,27777-8212,2021-01-29,One_Time,Fundraise,Charge,200.00,In_Memory_Of,Mr. Smoochers,,
Dana,danish,dana@pastry.com,17722 Fifth of Gin,,Spayallup,WA,97777-4132,2021-01-30,One_Time,Fundraise,Charge,35.00,In_Memory_Of,Killer Kitty,,
```

## Advanced usage

If you'll be using this app a lot, then it will be a good idea to create a native program.

```bash
go build -o ~/go/bin/fundraise_dedication cmd/activity/fundraise/dedication/main.go
```

The output will be an executable in `~/go/bin` in your home directory.
Add `~/go/bin` to the PATH list that your OS uses and you'll be able to invoke the program with a command like this.

```bash
fundraise_dedication --login company.yaml
```

## Questions?  Comments?

Use the [GitHub issues page](https://github.com/salsalabs/goengage/issues) to report problems, ask questions or make comments. Please don't bother the nice folks at Salsalabs.  This is their nesting season and they will bite intruders.
