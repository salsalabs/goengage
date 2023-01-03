#!/bin/bash

# Salsa has agreed to provide a list of dedicated (IHO/IMO) donations
# on a weekly basis. The donations are gathered into a CSV file by a
# Go application.  You can see the source for the application by
# navigating here:
#
#https://github.com/salsalabs/goengage/tree/master/cmd/activity/fundraise/dedication
#
# The CSV is renamed with the current date, ready for retrieval.
[ ! -e ~/tmp/apda ] && mkdir tmp/apda
cd ~/tmp/apda
export PATH=~/go/bin:$PATH
fundraise_dedication --login ~/.logins/apda.yaml $*
