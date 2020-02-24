#!/bin/bash

# -------------------------------------------------------------------
# Function to submit a request and display the total and count from
# the response.
#
# $1: API host
# $2: API token
# $3: API endpoint
# $4: request payload
# $5: (optional) command for jq
# -------------------------------------------------------------------
function go () {
    echo Request
    echo $4 | jq 
    echo Response
    curl -s -H "Content-Type: application/json" \
    -H "`echo authToken: $2`" \
    -X POST \
    -d "$4" \
    https://$1$3 \
    | jq "$5"
}

# -------------------------------------------------------------------
# Setup
# -------------------------------------------------------------------
token='wBTvk4rH5auTh4up8nOaVCcJBYWT3jr2Wk7QnlcOc4QnN1kWkb5hzf4Jge-_hHpaCEqwKmhH_Y953ExF9VdQ8MuR7dE3It-UwCBEnK4tUj8'
host='hq.uat.igniteaction.net'
endpoint='/api/integration/ext/v1/activities/search'
read -r -d '' request <<'EOF'
{
    "header": {
        "refId": ""
    },
    "payload": {
        "type": "FUNDRAISE",
        "offset": 0,
        "count": 20,
        "modifiedFrom": "2000-01-01T00:00:00.000Z"
    }
}
EOF
jq_command='.payload.activities[]| {actvitiyId: .activityId,formName: .activityFormName, transactions: .transactions[]|{date: .date, type: .type, amount: .amount}}'


go "$host" "$token" "$endpoint" "$request" "$jq_command"
