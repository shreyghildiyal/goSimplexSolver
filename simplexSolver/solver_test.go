package simplexsolver

import (
	"testing"
)

func TestSolver_GetSolution(t *testing.T) {

	constraints := []Equation{
		{
			rhs: 6,
			lhs: map[string]float64{
				"x": 3,
				"y": 1,
			},
			comparator: LTE,
		},
		{
			rhs: 7,
			lhs: map[string]float64{
				"x": 1,
				"y": 2,
			},
			comparator: LTE,
		},
	}

	objectiveFunc := map[string]float64{
		"x": 2,
		"y": 1,
	}

	expectedMax := 5.0
	expectedDistribution := map[string]float64{
		"x": 1.0,
		"y": 3.0,
	}

	s := SimplexSolver{constraints: constraints, objectiveFunction: objectiveFunc}

	max, vals := s.GetSolution2()

	if max != expectedMax {
		t.Error("max was incorrect")
	}
	for k, v := range vals {
		if ev, ok := expectedDistribution[k]; ok {
			if ev != v {
				t.Errorf("The value for %s is incorrect. Expected %f, found %f", k, expectedDistribution[k], v)
			}
		} else {
			t.Errorf("Resultant distribution contains variable %s but it is not present in expected distribution", k)
		}
	}
}

func TestSolver_GetSolution2(t *testing.T) {

	constraints := []Equation{
		{
			rhs: 15,
			lhs: map[string]float64{
				"x": 1,
				"y": 3,
			},
			comparator: LTE,
		},
		{
			rhs: 28,
			lhs: map[string]float64{
				"x": 2,
				"y": 5,
			},
			comparator: LTE,
		},
	}

	objectiveFunc := map[string]float64{
		"x": 1,
		"y": 4,
	}

	expectedMax := 20.0
	expectedDistribution := map[string]float64{
		"x": 0.0,
		"y": 5.0,
	}

	s := SimplexSolver{constraints: constraints, objectiveFunction: objectiveFunc}

	max, vals := s.GetSolution2()

	if max != expectedMax {
		t.Error("max was incorrect")
	}
	for k, v := range vals {
		if ev, ok := expectedDistribution[k]; ok {
			if ev != v {
				t.Errorf("The value for %s is incorrect. Expected %f, found %f", k, expectedDistribution[k], v)
			}
		} else {
			t.Errorf("Resultant distribution contains variable %s but it is not present in expected distribution", k)
		}
	}
}
