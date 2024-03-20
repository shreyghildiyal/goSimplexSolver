package simplexsolver

type Comparison int

const (
	LT = iota
	LTE
	EQ
	GTE
	GT
)

func (c Comparison) String() string {
	return [...]string{"Less Than", "LessThanEqualTo", "EqualTo", "GreaterThanEqualTo", "GreaterThan"}[c-1]
}

type Equation struct {
	lhs        map[string]float64
	inequality Comparison
	rhs        float64
}

func (eq Equation) GetLhs() map[string]float64 {
	return eq.lhs
}

func (eq Equation) GetInequality() Comparison {
	return eq.inequality
}
func (eq Equation) GetRhs() float64 {
	return eq.rhs
}
