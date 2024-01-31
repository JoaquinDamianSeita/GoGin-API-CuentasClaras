package services

import (
	"time"
)

const INCOME_TYPE string = "income"
const EXPENSE_TYPE string = "expense"

var utcLocation, _ = time.LoadLocation("UTC")
