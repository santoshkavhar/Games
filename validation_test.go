package main

import(
	"testing"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano() ) ) 

func TestCompleteValidationValidGrid(t *testing.T){
	
	s := createSudoku()

	if CompleteValidation(&s.grid) != true{
		t.Error("Valid Grid,  function should return true!")
	}
}

func TestCompleteValidationInvalidGrid(t *testing.T){
	
	s := createInvalidSudoku()

	if CompleteValidation(&s.grid) != false{
		t.Error("Invalid Grid, function should return false!")
	}
}

func TestCheckForViolationValidInputFromUser(t *testing.T){
	
	s := createValidRowColumnBlockSudoku()
	
	if checkForViolation(&s.grid, 0, 0) != true{
		t.Error("Valid UserInput,  function should return true!")
	}	
}

func TestCheckForViolationInvalidInputFromUser(t *testing.T){
	
	s := createInvalidRowColumnBlockSudoku()
	
	if checkForViolation(&s.grid, 0, 0) != false{
		t.Error("Invalid UserInput,  function should return false!")
	}	
}

func TestValidateCorrectRows(t *testing.T){
	
	s := createSudoku()
	
	if validateRows(&s.grid) != true{
		t.Error("All Rows were Valid,  function should return true!")
	}	
}

func TestValidateIncorrectRows(t *testing.T){
	
	s := createInvalidSudoku()
	
	if validateRows(&s.grid) != false{
		t.Error("Atleast one of the Rows were Invalid,  function should return false!")
	}	
}

func TestValidateCorrectColumns(t *testing.T){
	
	s := createSudoku()
	
	if validateColumns(&s.grid) != true{
		t.Error("All Columns were Valid,  function should return true!")
	}	
}

func TestValidateIncorrectColumns(t *testing.T){
	
	s := createInvalidSudoku()
	
	if validateColumns(&s.grid) != false{
		t.Error("Atleast one of the Columns were Invalid,  function should return false!")
	}	
}

func TestValidateCorrectBlocks(t *testing.T){
	
	s := createSudoku()
	
	if validateBlocks(&s.grid) != true{
		t.Error("All Blocks were Valid,  function should return true!")
	}	
}

func TestValidateIncorrectBlocks(t *testing.T){
	
	s := createInvalidSudoku()
	
	if validateBlocks(&s.grid) != false{
		t.Error("Atleast one of the Blocks were Invalid,  function should return false!")
	}	
}

func TestValidateCorrectSingleRow(t *testing.T){
	
	s := createSudoku()
	
	// Passing Random RowNum and whatType = 1
	if validateSingleRow(&s.grid,r.Intn(9), 1 ) != true{
		t.Error("The Row was Valid,  function should return true!")
	}	
}

func TestValidateIncorrectSingleRow(t *testing.T){
	// Column Of Zeroes so that none of the rows is valid
	s := createColumnOfRandomsSudoku()
	
	// Passing Random RowNum and whatType = 1
	if validateSingleRow(&s.grid,r.Intn(9), 1 ) != false{
		t.Error("The Row was Invalid,  function should return false!")
	}	
}

func TestValidateCorrectSingleColumn(t *testing.T){
	
	s := createSudoku()
	
	// Passing Random ColumnNum and whatType = 1
	if validateSingleColumn(&s.grid,r.Intn(9), 1 ) != true{
		t.Error("The Column was Valid,  function should return true!")
	}	
}

func TestValidateIncorrectSingleColumn(t *testing.T){
	// Row Of Randoms so that none of the Columns is valid
	s := createRowOfRandomsSudoku()
	
	// Passing Random ColumnNum and whatType = 1
	if validateSingleColumn(&s.grid,r.Intn(9), 1 ) != false{
		t.Error("The Column was Invalid,  function should return false!")
	}	
}

func TestValidateCorrectSingleBlock(t *testing.T){
	
	s := createSudoku()
	
	// Passing Random BlockNum Offset for posX, posY and whatType = 1
	if validateSingleBlock(&s.grid,r.Intn(3) * 3, r.Intn(3) * 3 , 1 ) != true{
		t.Error("Every Block was Valid,  function should return true!")
	}	
}

func TestValidateIncorrectSingleBlock(t *testing.T){
	// 
	s := createBlockOfRandomsSudoku()
	
	// Passing Random BlockNum Offset for posX, posY and whatType = 1
	if validateSingleBlock(&s.grid,r.Intn(3) * 3, r.Intn(3) * 3 , 1 ) != false{
		t.Error("Every Block was Invalid,  function should return false!")
	}	
}









func createSudoku() ( *Sudoku){
	s := Sudoku{}
	s.initializeSudoku()
	s.generateGrid()
	
	return &s
}

func createInvalidSudoku() ( *Sudoku){
	s := createSudoku()
	// Creating a Random location as Invalid entry
	row := r.Intn(9)
	col := r.Intn(9)
	s.grid[row][col] = 0	
	
	return s
}


func createValidRowColumnBlockSudoku() ( *Sudoku){
	s := createSudoku()
	
	// Recreating the 0th Row
	for i:=0; i < 9; i++ {
		s.grid[0][i] = i+1
	}
	// Recreating the 0th Column, leaving grid[0][0] entry
	for i:= 1; i < 9; i++ {
		s.grid[i][0] = 10 - i 
	}
	// Recreating the entries grid[1][1], grid[1][2], grid[2][1], grid[2][2]
	value := 4
	for i:= 1; i< 3; i++{
		for j:= 1; j< 3; j, value = j+1, value+1{
			s.grid[i][j] = value 
		}
	}
	

	/*
	// Our Grid will have a similar look
	
	1	9	8	7	6	5	4	3	2	
	2	4	5
	3	6	7
	
	4
	5
	6
	
	7
	8
	9
	*/ 
	return s
}



func createInvalidRowColumnBlockSudoku() ( *Sudoku){
	s := createSudoku()
	
	// Recreating the 0th Row
	for i:=0; i < 9; i++ {
		s.grid[0][i] = i+1
	}
	// Below Code Makes the block invalid
	// Recreating the 0th Column, leaving grid[0][0] entry
	for i:= 0; i < 9; i++ {
		s.grid[i][0] = i + 1 
	}
	// Recreating the entries grid[1][1], grid[1][2], grid[2][1], grid[2][2]
	value := 4
	for i:= 1; i< 3; i++{
		for j:= 1; j< 3; j, value = j+1, value+1{
			s.grid[i][j] = value 
		}
	}
	

	/*
	// Our Grid will have a similar look
	
	1	2	3	7	6	5	4	3	2	
	2	4	5
	3	6	7
	
	4
	5
	6
	
	7
	8
	9
	*/ 
	return s
}

func createColumnOfRandomsSudoku() ( *Sudoku){
	s := createSudoku()
	
	// Recreating a Random Column any Random Values >= 10
	col := r.Intn(9)
	for i:=0; i < 9; i++ {
		s.grid[i][col] = r.Intn(9) + 10
	}
	
	return s	
}

func createRowOfRandomsSudoku() ( *Sudoku){
	s := createSudoku()
	
	// Recreating a Random Row with any Random Values >= 10
	row := r.Intn(9)
	for i:=0; i < 9; i++ {
		s.grid[row][i] = r.Intn(9) + 10
	}
	
	return s	
}


func createBlockOfRandomsSudoku() ( *Sudoku){
	s := createSudoku()
	
	// Recreating a Grid with any Random element of each block 
	// being any invalid entry with Value >= 10
	for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				s.grid[i*3 + r.Intn(3)][j*3 + r.Intn(3) ] = r.Intn(9) + 10
			}
	}
	
	return s	
}



