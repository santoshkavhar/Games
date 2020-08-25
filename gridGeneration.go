package main

import (	
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// Initialization Logic

func (s *Sudoku) initializeSudoku() {	
	for i := 0; i < 9; i++ {
		s.oneToNineArray[i] = i + 1
		// initializing the grid and UserGrid to all zeroes
		for j := 0; j < 9; j++ {
			s.grid[i][j] = 0
			s.userGrid[i][j] = 0
		}
	}
}
	

func(s *Sudoku) generateGrid() error {
	// WaitGroup is used as we will be using go routines in Grid Generation
	var wg sync.WaitGroup

	// seedValue is what we pass to random function for first three block generation
	// r is used to get a random Sudoku from by selecting a random Seed from valid Seeds
	//  r.Intn(N) will return any number in range [0,N)
	r := rand.New(rand.NewSource(time.Now().UnixNano() ) )
	
	// randomValue should range from 0  to len(seeds) - 1 only
	randomValue :=r.Intn( len(seeds) - 1 )	

	seedValue := seeds[randomValue]

	/*
	   Consider our Grid consists of 9 blocks, such as

	   block0	block1		block2
	   block3	block4		block5
	   block6	block7		block8

	   Each block consists of a 3 X 3 Array of elements from 1 to 9, with no duplicates allowed

	*/
	
	// Below logic is for block number 0, 4, 8
	
	// for loop is run 3 times as we will be initializing 3 blocks at first

	
	// Logic which fills block 0, 4, 8 of the suduko
	
	// offset will help us identify block number 0, 4, 8
	// since we are generating blocks independently, we can use go routines
	for offset := 0; offset <= 6; offset += 3 {
		// Add 1 process to waitGroup named wg
		wg.Add(1)
		// Pass WaitGroup's address as that will be modified by thread
		// i.e It will get decremented after thread execution completes
		go s.parallelBlock048(offset, seedValue, &wg)
	}

	// wait for 3 threads to finish execution
	wg.Wait()

	
	// Logic for generating block 2, 3, 7 of the suduko
	
	blockID := 2
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)

	blockID = 3
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)

	blockID = 7
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)

	wg.Wait()


	// Logic for generating block 1, 5, 6

	blockID = 1
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)

	blockID = 5
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)

	blockID = 6
	wg.Add(1)
	go s.parallelBlock237(blockID, &wg)
	
	wg.Wait()

	// generateGrid might not generate a valid grid without a valid seed
	// So, just as a precaution Error is returned
	// CompleteValidation Method will return false if Sudoku is invalid
	// Passing address of grid, so as to avoid overhead of allocation of new grid
	sudukoGenerated := CompleteValidation(&s.grid)
	if sudukoGenerated != true {
		fmt.Println(seedValue)
		displayGrid(&s.grid)
		return fmt.Errorf("Sudoku wasn't generated!")		
	}


	return nil
}


func (s *Sudoku) getGridForUser(){
	// copying grid to userGrid
	s.userGrid = s.grid
	r := rand.New(rand.NewSource(time.Now().UnixNano() ) ) 
	
	// For Easy level, upto 5 boxes will need filling
	// For Medium Level upto 20 boxes will need filling
	// For Hard level upto 35 boxes will need filling
	
	// Error Handling
	if s.level > 2 {
		// changing level to hard only
		s.level = 2
	}
	blankBoxes := s.level * 15 + 5
	
	for  blankBoxes > 0 {
		row := r.Intn(9)
		col := r.Intn(9)
		// decrement blankBoxes only when suitable entry is found
		// else regenerate the randoms
		if s.userGrid[row][col] != 0 {
			s.userGrid[row][col] = 0
			blankBoxes--
			
		}
	}
}


func getStringArray(userGrid [9][9]int) string {
	var str string
	for i :=0 ; i <9 ; i++ {
		for j :=0; j<9; j++ {
			str = str + strconv.Itoa(userGrid[i][j])
		}
	}
	return str
}

// Check if user has won
func (s *Sudoku) checkWin() bool {

	// If userGrid is successfully validated, return true
	return CompleteValidation(&s.userGrid)
}


// Printing the grid
func displayGrid(grid *[9][9]int) {

	/*
			// for SeedValue = 560
	        1       8       3       6       7       9       4       2       5
	        9       5       2       8       1       4       3       7       6
	        4       7       6       3       5       2       8       9       1
	        2       1       7       4       3       5       9       6       8
	        6       9       8       1       2       7       5       3       4
	        3       4       5       9       8       6       7       1       2
	        8       6       9       2       4       3       1       5       7
	        7       3       4       5       6       1       2       8       9
	        5       2       1       7       9       8       6       4       3
	*/
	
	
	fmt.Println()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print("\t", grid[i][j])
		}
		fmt.Println()
	}
}


func  (s *Sudoku) parallelBlock048(offset, seedValue int, wg *sync.WaitGroup) {
	
	// b is copy of array
	b := s.oneToNineArray 
	
	// c is a slice referring array b
	c := b[:]           
	
	rand.Seed(int64(offset + seedValue))
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
		
			// getting a random index i
			randomIndex := rand.Intn(len(c))
			
			// Putting the element at i into grid
			// Putting the element onto grid
			s.grid[i+offset][j+offset] = c[randomIndex]
			
			// Put last element at index i
			c[randomIndex] = c[len(c)-1]
			
			// Decrement size of slice by 1
			c = c[:len(c)-1]
		}
	}
	wg.Done()
	return
}

// To Do -> Make code Cleaner
func  (s *Sudoku) parallelBlock237(blockID int, wg *sync.WaitGroup) {

	// posX will contain index of the first element of the block with id as blockID
	// Same is true for posY
	posX, posY := (blockID/3)*3, (blockID%3)*3
	
	// creating map for each block entry
	mapp := make(map[int][]int)
	//fmt.Println(mapp)
	
	// vg is used for waiting on entering Row and Column entries into mapp

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			/*
			validation logic
			mapIndex is where our entries of each Block get IDs from 0 to 8
			0	1	2
			3	4	5
			6	7	8
			*/
			mapIndex := i*3 + j
			// for Row and Cloumn
			s.enterRow(posX+i, mapIndex, mapp )
			s.enterColumn(posY+j, mapIndex, mapp)
		}
	}
	
	// Copy each map to remove duplicates
	for i := 0; i < 9; i++ {
		map2 := make(map[int]bool)
		val, ok := mapp[i]
		if ok {
			for j := 0; j < len(val); j++ {
				map2[val[j]] = true
			}
			delete(mapp, i)
			for key, _ := range map2 {
				mapp[i] = append(mapp[i], key)
			}
		} else {
			// No entry in mapp[i] means something went wrong
			// This won't happen as seeds are predefined
			// Below return is just for future testing
			return
		}
	}

	// inserting values into Suduko
	s.insert(blockID, mapp)

	wg.Done()
	return
}

func  (s *Sudoku) insert(blockID int, mapp map[int][]int) bool {

	arrayOfSlices := [9][]int{}
	sequence := [9]int{}
	// create array of slices which shows, which elements ought to be present in any block
	// loop for 9 times
	// if the element is present in the mapp, then just avoid adding it to arrayOfSlices

	for i, vals := range mapp {
		// beacuse of 1 to 9 index, avoiding any complexity
		auxArray := [10]int{}
		for j := 0; j < len(vals); j++ {
			auxArray[vals[j]] = 1
			// lesson learnt : avoid using range for slices
		}
		for k := 1; k < 10; k++ {
			if auxArray[k] == 0 {
				arrayOfSlices[i] = append(arrayOfSlices[i], k)
			}
		}
	}

	//fmt.Println(arrayOfSlices)
	// initializing
	for i := 0; i < 9; i++ {
		if len(arrayOfSlices[i]) > 0 {
			sequence[i] = arrayOfSlices[i][0]
		} else {
			return false
		}
	}
	//fmt.Println("Initial Sequence ", sequence)

	// using backtracking
	//fmt.Println(sequence)
	value := solveArray(&sequence, arrayOfSlices, 0, blockID)
	if value == false {
		//fmt.Println("No solution!")
	} else {
		//fmt.Println("Sequence :", sequence)
	}

	// Adding elements to the Grid

	posX, posY := (blockID/3)*3, (blockID%3)*3
	//b:= oneToNineArray[:]
	//counter := [9]int{}
	k := 0
	for i := posX; i < posX+3; i++ {
		for j := posY; j < posY+3; j++ {
			// copy 3 X 3 array to grid
			s.grid[i][j] = sequence[k]
			k++
		}
	}
	return true
}

// This function is used for backtracking entries within a block only

// Logic used:
// Candidate entries are array of Slices
// We simply create a combination of 9 of these entries, if a combination is safe
// we return true, if neither was valid then return false

//This combination is stored in variable named sequence
func solveArray(sequence *[9]int, arrayOfSlices [9][]int, k, blockID int) bool {
	if k >= 9 {
		return true
	}
	for i := 0; i < len(arrayOfSlices[k]); i++ {
		// move ahead only if entry is safe
		if isSafe(sequence, arrayOfSlices, k, i, blockID) {
			sequence[k] = arrayOfSlices[k][i]
			//fmt.Println(sequence[k])

			if solveArray(sequence, arrayOfSlices, k+1, blockID) {
				return true
			}
			//sequence[k] = arrayOfSlices[k][i+1]
		}
	}
	return false
}

// This is a helper function to solveArray Function 
func isSafe(sequence *[9]int, arrayOfSlices [9][]int, k, index, blockID int) bool {
	if k == 0 {
		return true
	}
	// copySequence is also the same array
	copySequence := sequence[:k]
	
	// Checking if sequence will be safe after changing a entry
	copySequence = append(copySequence, arrayOfSlices[k][index])

	mapp := make(map[int]bool)
	for i := 0; i < len(copySequence); i++ {
		mapp[copySequence[i]] = true
	}

	if len(mapp) == k+1 {
		return true
	}
	return false
}

func (s *Sudoku)  enterColumn(ColumnNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		if s.grid[i][ColumnNum] != 0 {
			mapp[mapIndex] = append(mapp[mapIndex], s.grid[i][ColumnNum])
		}
	}
}

func (s *Sudoku)  enterRow(RowNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		if s.grid[RowNum][i] != 0 {
			mapp[mapIndex] = append(mapp[mapIndex], s.grid[RowNum][i])
		}
	}
}
