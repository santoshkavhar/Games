package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
)

/*
// offsets for the blocks
const (
	blockID1 = iota * 3	// 0
	blockID5		// 3
	blockID9		// 6
)
*/

var (
	oneToNineArray = [9]int{}
	grid           = [9][9]int{}
)

func main() {

	//
	var wg sync.WaitGroup

	//
	// Initialization Logic
	//
	for i := 0; i < 9; i++ {
		oneToNineArray[i] = i + 1
		// initializing the matrix to blank := 0
		for j := 0; j < 9; j++ {
			grid[i][j] = 0
		}
	}

	/*
	   Consider our Grid consists of 9 blocks, such as

	   block0	block1		block2
	   block3	block4		block5
	   block6	block7		block8

	   Each block consists of a 3 X 3 Array of elements from 1 to 9, with no duplicates allowed

	*/
	// Our index range is between 0 and len(c)
	// Below logic is for block number 0 , 4, 8

	// for loop is run 3 times as we will be initializing 3 blocks at first

	//
	// Logic which fills block 0, 4, 8 of the suduko
	//
	for offset := 0; offset <= 6; offset += 3 {
		wg.Add(1)
		go parallelBlock048(offset, &wg)
	}

	// wait for the threads
	wg.Wait()

	//
	// Logic which fills block 2, 3, 7 of the suduko
	//
	blockID := 2
	wg.Add(1)
	go parallelBlock237(blockID, &wg)
	
	blockID = 3
	wg.Add(1)
	go parallelBlock237(blockID, &wg)

	blockID = 7
	wg.Add(1)
	go parallelBlock237(blockID, &wg)
	

	wg.Wait()

	/*
		1       3       9       0       0       0       0       0       0
		5       7       6       0       0       0       0       0       0
		2       8       4       0       0       0       0       0       0
		0       0       0       5       2       3       0       0       0
		0       0       0       1       7       4       0       0       0
		0       0       0       6       8       9       0       0       0
		0       0       0       0       0       0       7       4       3
		0       0       0       0       0       0       5       9       8
		0       0       0       0       0       0       6       1       2
	*/
	// Printing the grid

	fmt.Println()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print("\t", grid[i][j])
		}
		fmt.Println()
	}

	// Validation logic
	sudukoValidation := CompleteValidation()

	fmt.Println(sudukoValidation)
}

func parallelBlock048(offset int, wg *sync.WaitGroup) {

	b := oneToNineArray // b is copy of array
	c := b[:]           // c is a slice referring array b
	rand.Seed(int64(offset))
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// getting a random index i
			randomIndex := rand.Intn(len(c))
			// Putting the element at i into grid
			// Putting the element onto grid
			grid[i+offset][j+offset] = c[randomIndex]
			// Put last element at index i
			c[randomIndex] = c[len(c)-1]
			// Decrement size of slice by 1
			c = c[:len(c)-1]
		}
	}
	wg.Done()
	return
}

// needs editing
func parallelBlock237(blockID int, wg *sync.WaitGroup) {

	// creating map for each block entry

	posX, posY := (blockID/3)*3, (blockID%3)*3
	mapp := make(map[int][]int)
	//fmt.Println(mapp)

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
			enterRow(posX+i, mapIndex, mapp)
			enterColumn(posY+j, mapIndex, mapp)
			// For Block
			/*if grid[posX+i][posY+j] != 0{
				for  k :=0 ; k < 9 ; k++{
					if mapIndex != k {
						mapp[mapIndex] = append(mapp[mapIndex] , grid[posX+i][posY+j] )
					}
				}
			}
			*/
			//enterBlock(posX + i, posY+ j, blockID, mapIndex ,  mapp)
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
		}
	}

	// putting values into Suduko

	value := insert(blockID, mapp)
	if value == false {
		fmt.Println("Suduko not possible, Try with other seed")
		os.Exit(-1)
	}

	//slices := [][]entries
	//fmt.Println(mapp)

	/*
		b := oneToNineArray // b is copy of array
		c := b[:]           // c is a slice referring array b

		int i := 0

		//slices := [][]int{}
	*/
	wg.Done()
	return
}

func insert(blockID int, mapp map[int][]int) bool {

	arrayOfSlices := [9][]int{}
	var sequence = [9]int{}
	// create array of slices which shows, which elements ought to be present
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
	// initializing
	for i := 0; i < 9; i++ {
		sequence[i] = arrayOfSlices[i][0]
	}
	fmt.Println("Initial Sequence ", sequence)

	// using backtracking
	//fmt.Println(sequence)
	value := solveArray(&sequence, arrayOfSlices, 0)
	if value == false {
		fmt.Println("No solution!")
	} else {
		fmt.Println("Sequence :", sequence)
	}

	// Adding elements

	posX, posY := (blockID/3)*3, (blockID%3)*3
	//b:= oneToNineArray[:]
	//counter := [9]int{}
	k := 0
	for i := posX; i < posX+3; i++ {
		for j := posY; j < posY+3; j++ {
			// copy 3 X 3 array to grid
			grid[i][j] = sequence[k]
			k++
		}
	}
	return true
}

func solveArray(sequence *[9]int, arrayOfSlices [9][]int, k int) bool {
	if k >= 9 {
		return true
	}
	for i := 0; i < len(arrayOfSlices[k]); i++ {
		// move ahead only if entry is safe
		if isSafe(sequence, arrayOfSlices, k, i) {
			sequence[k] = arrayOfSlices[k][i]
			//fmt.Println(sequence[k])

			if solveArray(sequence, arrayOfSlices, k+1) {
				return true
			}
			//sequence[k] = arrayOfSlices[k][i+1]
		}
	}
	return false
}

func isSafe(sequence *[9]int, arrayOfSlices [9][]int, k, index int) bool {
	if k == 0 {
		return true
	}
	// copySequence is also the same array
	copySequence := sequence[:k]
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

/*
	index := 0
			k := 0
			for  {
				slice := mapp[index]
				found := false
				for l:=0 ; l <len(slice); l++{
					if b[k] == slice[l] {
						found = true
						k++
						break
					}
				}
				if !found {
					grid[i][j] = b[k]
					counter[index] = k
					index++
					b[k] = b[len(b) - 1]
					b = b[:len(b) - 1]
					break
				}
			}


*/

/*
func checkSafe()
*/

/*
func enterRow(RowNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		// if element is 0 , avoid adding it to the mapp
		if grid[RowNum][i] == 0 {
			continue
		}


		if !found {
			mapp[mapIndex] = append(mapp[mapIndex], grid[RowNum][i])
		}
	}
}
*/
/*
func enterColumn(ColumnNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		if grid[i][Column] == 0 {
			continue
		}
		found := false
		val, ok := mapp[mapIndex]
		if ok {
			for j := 0; j < len(val); j++ {
				if grid[i][Column] == val[j] {
					found = true
					break
				}
			}
		}

		if !found {
			mapp[mapIndex] = append(mapp[mapIndex], grid[i][Column])
		}

	}
}
*/

func enterColumn(ColumnNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		if grid[i][ColumnNum] != 0 {
			mapp[mapIndex] = append(mapp[mapIndex], grid[i][ColumnNum])
		}
	}
}

func enterRow(RowNum, mapIndex int, mapp map[int][]int) {
	for i := 0; i < 9; i++ {
		if grid[RowNum][i] != 0 {
			mapp[mapIndex] = append(mapp[mapIndex], grid[RowNum][i])
		}
	}
}

/*
func enterBlock(posX, posY , blockID , mapIndex int, mapp map[int][]int) {

	/*for i := posX; i < posX + 3; i++ {
		for j :=posY ; j <posY + 3 ; j++ {
			if grid[posX][posY] != 0 {
				mapp[mapIndex] = append(mapp[mapIndex], grid[posX][posY])
			}
		}
	}
}


func enterBlock(posX, posY int, mapp map[int][]int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {

			// Copy Block entry to each map
			if grid[i+posX][j+posY] == 0 {
				continue
			}
			mapIndex := i*3 + j

			for k := 0; k < 3; k++ {

				if mapIndex != k  && !found{
					mapp[k] = append(mapp[k], grid[i+posX][j+posY])
				}
			}

			found := false
			val, ok := mapp[mapIndex]
			if ok {
				for k := 0; k < len(val); k++ {
					if grid[i+posX][j+posY] == val[k] {
						found = true
						break
					}
				}
			}


		}
	}
}
*/

// Validation Logic
func CompleteValidation() bool {

	//truthValue := true

	truthValue := validateRows() && validateColumns() && validateBlocks()

	return truthValue

}

func validateRows() bool {
	for i := 0; i < 9; i++ {
		if truthValueRow := validateSingleRow(i); truthValueRow == false {
			return false
		}
	}
	return true
}

func validateSingleRow(RowNum int) bool {
	validationMap := make(map[int]bool, 9)
	for j := 0; j < 9; j++ {
		if grid[RowNum][j] >= 1 && grid[RowNum][j] <= 9 {
			validationMap[grid[RowNum][j]] = true
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
	for i := 0; i < 9; i++ {
		if truthValueColumn := validateSingleColumn(i); truthValueColumn == false {
			return false
		}
	}
	return true
}

func validateSingleColumn(ColumnNum int) bool {

	validationMap := make(map[int]bool, 9)
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
	return false
}

func validateBlocks() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Multilying by 3 , so as to get proper offset
			if truthValueBlock := validateSingleBlock(i*3, j*3); truthValueBlock == false {
				return false
			}
		}
	}
	return true
}

func validateSingleBlock(offsetX, offsetY int) bool {

	validationMap := make(map[int]bool, 9)
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
	return false
}
