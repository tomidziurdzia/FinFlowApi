package domain

type CategoryType int

const (
	CategoryTypeExpense CategoryType = iota
	CategoryTypeIncome
	CategoryTypeInvestment
)

func (ct CategoryType) String() string {
	switch ct {
	case CategoryTypeExpense:
		return "Expense"
	case CategoryTypeIncome:
		return "Income"
	case CategoryTypeInvestment:
		return "Investment"
	default:
		return "Unknown"
	}
}

func (ct CategoryType) Value() int {
	return int(ct)
}