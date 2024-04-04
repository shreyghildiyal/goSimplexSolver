package simplexsolver

type Comparison int

const (
	LTE Comparison = iota
	EQ
	GTE
)

func (c Comparison) String() string {
	return [...]string{"LessThanEqualTo", "EqualTo", "GreaterThanEqualTo"}[c-1]
}

type Equation struct {
	lhs        map[string]float64
	comparator Comparison
	rhs        float64
}

func (eq Equation) GetLhs() map[string]float64 {
	return eq.lhs
}

func (eq Equation) GetInequality() Comparison {
	return eq.comparator
}
func (eq Equation) GetRhs() float64 {
	return eq.rhs
}
