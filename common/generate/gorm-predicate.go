package generate

// GormPredicate is a string that acts as a condition in the where clause
type GormPredicate string

var (
	EqualPredicate              = GormPredicate("=")
	NotEqualPredicate           = GormPredicate("<>")
	GreaterThanPredicate        = GormPredicate(">")
	GreaterThanOrEqualPredicate = GormPredicate(">=")
	SmallerThanPredicate        = GormPredicate("<")
	SmallerThanOrEqualPredicate = GormPredicate("<=")
	LikePredicate               = GormPredicate("LIKE")
)
