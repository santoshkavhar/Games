// Validation logic
package main

func CompleteValidation() bool {
	
	truthValue := true

	truthValue = validateRows()  &&   validateColumns()  &&   validateBlocks()

	return truthValue

} 

func validateRows() bool {
	for i:=0 ; i < 9; i++ {
		if truthValueRow := validateSingleRow(i) ; truthValueRow == false {
			return false
		}
	}
	return true
}

func validateSingleRow(RowNum int) bool {
	validationMap := make(map[int]bool, 9)
	for j := 0; j< 9 ; j++ {
		if grid[RowNum][j] >= 1 && grid[RowNum][j] <=9 {
			validationMap[ grid[RowNum][j] ] = true
		} else {
			return false
		}
	}

	if len(validationMap) == 9 {
		return true
	} 
	return false
}

func validateColumns() bool {
	for i:=0 ; i < 9; i++ {
		if truthValueColumn := validateSingleColumn(i) ; truthValueColumn == false {
			return false
		}
	}
	return true
}

func validateSingleColumn(ColumnNum int) bool {

	validationMap := make(map[int]bool, 9)
	for j := 0; j< 9 ; j++ {
		if grid[j][ColumnNum] >= 1 && grid[j][ColumnNum] <=9 {
			validationMap[ grid[j][ColumnNum] ] = true
		} else {
			return false
		}
	}

	if len(validationMap) == 9 {
		return true
	} 
	return false
}


func validateBlocks() bool {
	for i:= 0; i< 3 ; i++ {
		for j :=0 ; j <3; j++ {
			// Multilying by 3 , so as to get proper offset
			if truthValueBlock := validateSingleBlock(i * 3, j * 3) ; truthValueBlock == false {
				return false
			}
		}
	}
	return true
}

func validateSingleBlock(offsetX, offsetY int) bool {

	validationMap := make(map[int]bool, 9)
	for i:= 0; i< 3 ; i++ {
		for j :=0 ; j <3; j++ {
			if grid[offsetX + i][offsetY + j] >= 1 && grid[offsetX + i][offsetY + j]  <=9 {
				validationMap[ grid[offsetX + i][offsetY + j]  ] = true
			} else {
				return false
			}
		}
	}
	if len(validationMap) == 9 {
		return true
	} 
	return false
}

