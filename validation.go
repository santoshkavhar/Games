package main

// This file contains only Validation logic 


// To use -> go routines
func CompleteValidation(grid *[9][9]int) bool {

	// grid is a pointer to a 2D 9 X 9 array.
	// It could either be a userGrid or our GameGrid(alias grid) 
	truthValue := validateRows(grid) && validateColumns(grid) && validateBlocks(grid)

	return truthValue
}

// For Ad-hoc validations of UserGrid
func checkForViolation(grid *[9][9]int, row, col int) bool {
	truthValue := validateSingleRow(grid, row, 1) && validateSingleColumn(grid, col, 1) &&
			 validateSingleBlock(grid, (row/3 ) * 3, (col/3 ) * 3, 1)
	// (row/3 ) * 3  is used so that offset is maintained properly
	
	return truthValue
}

func validateRows(grid *[9][9]int) bool {
	// Here i = 0 means 0th row validation
	for i := 0; i < 9; i++ {
		if truthValueRow := validateSingleRow(grid, i, 0); truthValueRow == false {
			return false
		}
	}
	return true
}

// whatType is used to know if we are validating only a part the complete grid
// whatType = 0 means normal grid, complete grid validation is being performed
// whatType = 1 means spontaneous(Ad-hoc) validation
func validateSingleRow(grid *[9][9]int, RowNum int, whatType int) bool {

	validationMap := make(map[int]bool, 9)

	switch (whatType){
	// Normal Validation, i.e Part Of Complete Validation
	// We will add each of the Row Element(between 1 to 9 only) to a map
	// If at the end length of map == 9, then Row Formation is valid
	// Logic useful only for grid and when whatType = 0
	case 0:	
		for j := 0; j < 9; j++ {
			if grid[RowNum][j] >= 1 && grid[RowNum][j] <= 9 {
				validationMap[ grid[RowNum][j] ] = true
			} else {
				return false
			}
		}

		if len(validationMap) == 9 {
			return true
		}
	// Ad-hoc Validation for single Row
	case 1:
		for j := 0; j < 9; j++ {
			if grid[RowNum][j] >= 1 && grid[RowNum][j] <= 9 {
				if _, ok := validationMap[ grid[RowNum][j] ] ; ok {
					// duplicate found!
					return false
				}
				// else , Add entry to Map
				validationMap[ grid[RowNum][j] ] = true
			}
		}
		return true	
	}
	
	return false
}

func validateColumns( grid *[9][9]int) bool {
	// Here i = 0 means 0th column validation
	for i := 0; i < 9; i++ {
		if truthValueColumn := validateSingleColumn(grid, i, 0); truthValueColumn == false {
			return false
		}
	}
	return true
}

func validateSingleColumn(grid *[9][9]int, ColumnNum int, whatType int) bool {
	validationMap := make(map[int]bool, 9)

	switch (whatType){	
	case 0:
	// Normal Validation, i.e Part Of Complete Validation	
	// We will add each of the Column Element(between 1 to 9 only) to a map
	// If at the end length of map == 9, then Column Formation is valid
	// Logic useful only for grid and when whatType = 0
		for j := 0; j < 9; j++ {
			if grid[j][ColumnNum] >= 1 && grid[j][ColumnNum] <= 9 {
				validationMap[grid[j][ColumnNum]] = true
			} else {
				return false
			}
		}

		if len(validationMap) == 9 {
			return true
		}

	case 1:
	// Ad-hoc Validation for single Column
		for j := 0; j < 9; j++ {
			if grid[j][ColumnNum] >= 1 && grid[j][ColumnNum] <= 9 {
				if _, ok := validationMap[ grid[j][ColumnNum] ] ; ok {
					// duplicate found!
					return false
				}
				// else , Add entry to Map
				validationMap[ grid[j][ColumnNum] ] = true
			}
		}
		return true	
	}
	return false
}


func validateBlocks(grid *[9][9]int ) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Multilying by 3 , so as to get proper offset
			if truthValueBlock := validateSingleBlock(grid, i*3, j*3, 0); truthValueBlock == false {
				return false
			}
		}
	}
	return true
}

func validateSingleBlock(grid *[9][9]int, offsetX, offsetY, whatType int) bool {
	validationMap := make(map[int]bool, 9)

	switch (whatType){	
	case 0:
	// Normal Validation, i.e Part Of Complete Validation	
	// We will add each of the Block Element(between 1 to 9 only) to a map
	// If at the end length of map == 9, then Block Formation is valid
	// Logic useful only for grid and when whatType = 0
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if grid[offsetX+i][offsetY+j] >= 1 && grid[offsetX+i][offsetY+j] <= 9 {
					validationMap[grid[offsetX+i][offsetY+j]] = true
				} else {
					return false
				}
			}
		}
		if len(validationMap) == 9 {
			return true
		}

	case 1:
	// Ad-hoc Validation for single Block
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if grid[offsetX+i][offsetY+j] >= 1 && grid[offsetX+i][offsetY+j] <= 9 {
					if _, ok := validationMap[ grid[offsetX+i][offsetY+j] ] ; ok {
						// duplicate found!
						return false
					}
					// else , Add entry to Map
					validationMap[ grid[offsetX+i][offsetY+j] ] = true
				}
			}
		}
		return true	
	}
	return false
}



