package simplexsolver

// constraints always have the implicit inequality that all variables are positive
type Solver struct {
	constraints       []Equation         // constraint equations
	objectiveFunction map[string]float64 // the polynomia to maximize
}

func (s Solver) GetConstraints() []Equation {
	return s.constraints
}

func (s Solver) GetObjectiveFunction() map[string]float64 {
	return s.objectiveFunction
}

func (s *Solver) AddConstraint(eq Equation) {
	s.constraints = append(s.constraints, eq)
}

func (s *Solver) GetSolution() (float64, map[string]float64) {

	tableu := [][]float64{}
	variablePositionMap := map[string]int{}
	varLenCount := 0
	rhs := []float64{}

	// add all base variables to the tableu
	for _, eq := range s.constraints {
		row := make([]float64, varLenCount)
		for k, v := range eq.lhs {
			if pos, ok := variablePositionMap[k]; ok {
				row[pos] = v
			} else {
				variablePositionMap[k] = varLenCount
				row = append(row, v)
				varLenCount++

			}
		}
		tableu = append(tableu, row)
		rhs = append(rhs, eq.rhs)
	}

	// add all slack variables to the tableu
	for i, eq := range s.constraints {
		if eq.inequality == GTE {
			for j := 0; j < len(s.constraints); j++ {
				if i == j {
					tableu[j] = append(tableu[j], 1)

				} else {
					tableu[j] = append(tableu[j], 0)
				}
			}
		} else {
			//TODO: add handling for other constraints types
		}
	}

	objectiveRow := make([]float64, len(tableu[0]))

	for k, v := range variablePositionMap {
		objectiveRow[v] = s.objectiveFunction[k] * (-1)
	}

	var objectiveRhs float64 = 0

	maxVal, res := loopForSolution(tableu, rhs, objectiveRow, objectiveRhs)

	retMap := map[string]float64{}

	for k, index := range variablePositionMap {
		retMap[k] = res[index]
	}

	return maxVal, retMap
}

func loopForSolution(tableu [][]float64, rhs []float64, objectiveRow []float64, objectiveRhs float64) (float64, []float64) {
	// panic("unimplemented")

	for !maximizationDone(objectiveRow) {
		// choose pivot column
		pivotColumn := choosePivotColumn(objectiveRow)
		// choose pivot row
		pivotRow := choosePivotRow(tableu, pivotColumn, rhs)

		// make the pivot column zero in objectiveFunc
		objectiveRhs = objectiveRhs - rhs[pivotRow]*objectiveRow[pivotColumn]/tableu[pivotRow][pivotColumn]
		reduceObjectiveRow(objectiveRow, pivotColumn, tableu, pivotRow)

		// make pivot column in pivotRow equal to 1
		normalizePivotRow(tableu, pivotRow, pivotColumn)

		// make pivot column zero in all rows except pivotRow
		reduceNonPivotRows(tableu, pivotRow, pivotColumn)

	}

	basciColArr := []bool{}
	for i := 0; i < len(objectiveRow); i++ {
		basciColArr[i] = isColBasic(objectiveRow, tableu, i)
	}

	return 0, []float64{}
}

func isColBasic(objectiveRow []float64, tableu [][]float64, col int) bool {
	nonZeroCount := 0

	if objectiveRow[col] != 0 {
		nonZeroCount++
	}

	for row := 0; row < len(tableu); row++ {
		if tableu[row][col] != 0 {
			nonZeroCount++
			if nonZeroCount > 1 {
				return false
			}
		}
	}

	if nonZeroCount == 1 {
		return true
	} else {
		return false
	}

}

func choosePivotColumn(objectiveRow []float64) int {
	pivotColumn := 0
	minColVal := objectiveRow[0]
	for i, v := range objectiveRow {
		if v < minColVal {
			pivotColumn = i
			minColVal = v
		}
	}
	return pivotColumn
}

func choosePivotRow(tableu [][]float64, pivotColumn int, rhs []float64) int {
	pivotRow := -1
	var minRowVal float64 = -1
	foundOneRow := false

	for i := range tableu {
		if tableu[i][pivotColumn] > 0 && rhs[i] > 0 {
			if !foundOneRow {
				pivotRow = i
				minRowVal = rhs[i] / tableu[i][pivotColumn]
			} else {
				grad := rhs[i] / tableu[i][pivotColumn]
				if grad < minRowVal {
					pivotRow = i
					minRowVal = grad
				}
			}
		}
	}
	return pivotRow
}

func reduceObjectiveRow(objectiveRow []float64, pivotColumn int, tableu [][]float64, pivotRow int) {
	multiplier := objectiveRow[pivotColumn] / tableu[pivotRow][pivotColumn]
	for i := range objectiveRow {
		objectiveRow[i] = objectiveRow[i] - tableu[pivotRow][pivotColumn]*multiplier
	}
}

func normalizePivotRow(tableu [][]float64, pivotRow int, pivotColumn int) {
	pivotCellVal := tableu[pivotRow][pivotColumn]
	for col := 0; col < len(tableu[pivotRow]); col++ {
		tableu[pivotRow][col] = tableu[pivotRow][col] / pivotCellVal
	}
}

func reduceNonPivotRows(tableu [][]float64, pivotRow int, pivotColumn int) {
	for row := 0; row < len(tableu); row++ {

		if row != pivotRow {
			rowMultiplier := tableu[row][pivotColumn] / tableu[pivotRow][pivotColumn]
			for col := 0; col < len(tableu[row]); col++ {
				tableu[row][col] = tableu[row][col] - (rowMultiplier * tableu[pivotRow][col])
			}

		}

	}
}

func maximizationDone(objectiveRow []float64) bool {

	for _, val := range objectiveRow {
		if val < 0 {
			return false
		}
	}
	return true
}
