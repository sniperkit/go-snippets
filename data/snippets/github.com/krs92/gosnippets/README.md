# gosnippets
Go language  examples 


# 1) Identifying race conditions with race detection
>Assume that you start with $1,000 and attempt 200 $5 transactions. Each
transaction requires a query on the current balance of the account. If it passes, the
transaction is approved and $5 is removed from the balance. If it fails, the transaction is
declined and the balance remains unchanged.
we’ll build a bank account that starts with $1,000 and runs 100 transactions for a random amount between $0 and $25.


# 2) Dining philosophers problem

