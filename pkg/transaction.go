package goengage

//Describes search for transactions. Note that transactions are listed
//as part of actions *and* can be searched without actions being involved.

//Engage endpoints for transactions.
const (
	SearchTransactionDetails   = "/api/integration/ext/v1/transactionDetails/search"
	SearchTransactionTemplates = "/api/integration/ext/v1/transactionTemplates/search"
)
