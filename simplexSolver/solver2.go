package simplexsolver

func (s *SimplexSolver) GetSolution2() (float64, map[string]float64) {

	// find all the variables in constraints and objective function
	varPosMap := s.GetVarPosMap()
	// get tableu
	// 	row 0 -> artificial variable removal
	// 	row 1 -> objective function
	// row n -> constraints
	//		1 row per lte
	//		1 row per gte
	// 		2 rows per equality
	tableu, rhs := s.GetBasicTableu2(varPosMap)

	// do phase 1 of 2 phase simplex if needed
	if containsPositiveValue(tableu[0]) {
		tableu, rhs = minimizeArtificialVariables(tableu, rhs)
	}
	// do phase 2 of 2
	if containsNegativeValue(tableu[1]) {
		tableu, rhs = maximizeObjeciveFunction(tableu, rhs)
	}

	// interpret tableu
	maxValue := rhs[1]

	distribution := map[string]float64{}

	for varName, index := range varPosMap {
		if isBasic(tableu, index) {
			selectedRow := getNonZeroRow(tableu, index)
			distribution[varName] = rhs[selectedRow] / tableu[selectedRow][index]
		} else {
			distribution[varName] = 0
		}
	}

	return maxValue, distribution
}

func getNonZeroRow(tableu [][]float64, index int) int {
	for i := 1; i < len(tableu); i++ {
		if tableu[i][index] != 0 {
			return i
		}
	}
	panic("No non zero row present")
}

func isBasic(tableu [][]float64, index int) bool {
	nonZeroCount := 0
	for i := 1; i < len(tableu); i++ {
		if tableu[i][index] != 0 {
			nonZeroCount++
		}
	}

	return nonZeroCount == 1
}

func maximizeObjeciveFunction(tableu [][]float64, rhs []float64) ([][]float64, []float64) {

	for containsNegativeValue(tableu[1]) {
		pivotColumn := getNegPivotColumn(tableu[1])
		pivotRow := getPivotRow(tableu, rhs, pivotColumn)
		tableu, rhs = reduce(tableu, rhs, pivotRow, pivotColumn)
	}

	return tableu, rhs
}

func getNegPivotColumn(objectiuveRow []float64) int {
	pivotColumn := 0
	pivotVal := objectiuveRow[0]
	for i, val := range objectiuveRow {
		if val < pivotVal {
			pivotColumn = i
			pivotVal = val
		}
	}
	return pivotColumn
}

func minimizeArtificialVariables(tableu [][]float64, rhs []float64) ([][]float64, []float64) {

	for containsPositiveValue(tableu[0]) {
		pivotColumn := getPivotColumn(tableu[0])
		pivotRow := getPivotRow(tableu, rhs, pivotColumn)
		tableu, rhs = reduce(tableu, rhs, pivotRow, pivotColumn)
	}
	if rhs[0] > 0 {
		panic("Seems we cant get rid of the artificial variables. No solution might exist")
	}
	return tableu, rhs
}

func reduce(tableu [][]float64, rhs []float64, pivotRow, pivotColumn int) ([][]float64, []float64) {

	pivotVal := tableu[pivotRow][pivotColumn]
	// reduce the pivot row
	for col := range tableu[pivotRow] {
		tableu[pivotRow][col] = tableu[pivotRow][col] / pivotVal
	}
	rhs[pivotRow] = rhs[pivotRow] / pivotVal

	for row := range tableu {
		if row != pivotRow {
			if tableu[row][pivotColumn] != 0 {
				multiplier := tableu[row][pivotColumn] / tableu[pivotRow][pivotColumn]
				for col := range tableu[row] {
					tableu[row][col] = tableu[row][col] - multiplier*tableu[pivotRow][col]
				}
				rhs[row] = rhs[row] - multiplier*rhs[pivotRow]
			}
		}
	}

	return tableu, rhs
}

func getPivotRow(tableu [][]float64, rhs []float64, pivotColumn int) int {

	pivotRow := 2
	var pivotVal float64 = 0
	candidateFound := false
	for row := 2; row < len(tableu); row++ {
		if rhs[row] > 0 && tableu[row][pivotColumn] > 0 {
			v := rhs[row] / tableu[row][pivotColumn]
			if !candidateFound {
				pivotRow = row
				pivotVal = v
				candidateFound = true
			} else if v < pivotVal {
				pivotRow = row
				pivotVal = v
			}

		}
	}
	return pivotRow
}

func getPivotColumn(artificialRow []float64) int {
	pivotColumn := 0
	pivotVal := artificialRow[0]
	for i, val := range artificialRow {
		if val > pivotVal {
			pivotColumn = i
			pivotVal = val
		}
	}
	return pivotColumn
}

func containsPositiveValue(row []float64) bool {
	for _, v := range row {
		if v > 0 {
			return true
		}
	}
	return false
}

func containsNegativeValue(row []float64) bool {
	for _, v := range row {
		if v < 0 {
			return true
		}
	}
	return false
}

func (s *SimplexSolver) GetVarPosMap() map[string]int {
	varPosMap := map[string]int{}

	varCount := 0
	for varName := range s.objectiveFunction {
		if _, ok := varPosMap[varName]; !ok {
			varPosMap[varName] = varCount
			varCount++
		}
	}

	for _, constraint := range s.constraints {
		for varName := range constraint.lhs {
			if _, ok := varPosMap[varName]; !ok {
				varPosMap[varName] = varCount
				varCount++
			}
		}
	}

	return varPosMap
}

func (s *SimplexSolver) GetBasicTableu2(varPosMap map[string]int) ([][]float64, []float64) {
	tableu := [][]float64{}

	rhs := []float64{}

	// add zero row for artificaial constriant

	tableu, rhs = updateArtificalRow(varPosMap, tableu, rhs)

	// add row for objective function
	tableu, rhs = updateObjectiveRow(varPosMap, s.objectiveFunction, tableu, rhs)

	// start adding constraints

	for _, constraint := range s.constraints {

		if constraint.comparator != LTE && constraint.comparator != GTE && constraint.comparator != EQ {
			panic("unsupported comparator")
		}

		if constraint.comparator == LTE || constraint.comparator == EQ {
			tableu, rhs = updateLTEConstraint(tableu, constraint, varPosMap, rhs)

		}
		if constraint.comparator == GTE || constraint.comparator == EQ {

			tableu, rhs = updateGTEConstraint(tableu, constraint, varPosMap, rhs)
		}
	}

	return tableu, rhs
}

func updateGTEConstraint(tableu [][]float64, constraint Equation, varPosMap map[string]int, rhs []float64) ([][]float64, []float64) {
	row := make([]float64, len(tableu[0]))
	for k, v := range constraint.lhs {
		row[varPosMap[k]] = v
		tableu[0][varPosMap[k]] = tableu[0][varPosMap[k]] + v
	}
	tableu = append(tableu, row)
	rhs = append(rhs, constraint.rhs)
	rhs[0] = rhs[0] + constraint.rhs

	// add surplus variables
	for k := 0; k < len(tableu); k++ {
		if k == len(tableu)-1 || k == 0 {
			tableu[k] = append(tableu[k], -1)
		} else {
			tableu[k] = append(tableu[k], 0)
		}
	}

	// add artificial variables
	for k := 0; k < len(tableu); k++ {
		if k == len(tableu)-1 {
			tableu[k] = append(tableu[k], 1)
		} else {
			tableu[k] = append(tableu[k], 0)
		}
	}

	return tableu, rhs
}

func updateLTEConstraint(tableu [][]float64, constraint Equation, varPosMap map[string]int, rhs []float64) ([][]float64, []float64) {
	row := make([]float64, len(tableu[0]))
	for k, v := range constraint.lhs {
		row[varPosMap[k]] = v
	}
	tableu = append(tableu, row)

	// add slack variables
	for k := 0; k < len(tableu); k++ {
		if k == len(tableu)-1 {
			tableu[k] = append(tableu[k], 1)
		} else {
			tableu[k] = append(tableu[k], 0)
		}
	}
	rhs = append(rhs, constraint.rhs)
	return tableu, rhs
}

func updateObjectiveRow(varPosMap map[string]int, objectiveFunction map[string]float64, tableu [][]float64, rhs []float64) ([][]float64, []float64) {
	objectiveRow := make([]float64, len(varPosMap))
	for k, v := range objectiveFunction {
		pos := varPosMap[k]
		objectiveRow[pos] = v * -1
	}
	tableu = append(tableu, objectiveRow)
	rhs = append(rhs, 0)
	return tableu, rhs
}

func updateArtificalRow(varPosMap map[string]int, tableu [][]float64, rhs []float64) ([][]float64, []float64) {
	artificalRow := make([]float64, len(varPosMap))
	tableu = append(tableu, artificalRow)
	rhs = append(rhs, 0)
	return tableu, rhs
}
