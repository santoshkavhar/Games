package main

import (
	"fmt"
	"math/rand"
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
	grid = [9][9]int{}
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
	// Below logic is for block number 1, 5, 9

	// for is 3 times as we will be initializing 3 blocks at first

	//
	// Logic which fills block 0, 4, 8 of the suduko
	//
	for offset := 0; offset <= 6; offset += 3 {
		wg.Add(1)
		go parallelBlockFirst(offset, &wg)
	}

	// wait for the threads
	wg.Wait()


	//
	// Logic which fills block 2, 3, 7 of the suduko
	//
	/*blockID := 3
	wg.Add(1)
	go parallelBlockSecond(blockID / 3, blockID % 3, &wg)
	
	blockID = 4
	wg.Add(1)	
	go parallelBlockSecond(blockID / 3, blockID % 3, &wg)

	blockID = 8
	wg.Add(1)
	go parallelBlockSecond(blockID / 3, blockID % 3, &wg)
	
	
	wg.Wait()
	*/

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

func parallelBlockFirst(offset int, wg *sync.WaitGroup) {

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
/*
func parallelBlockSecond(posX , posY int, wg *sync.WaitGroup) {

	b := oneToNineArray // b is copy of array
	c := b[:]           // c is a slice referring array b
	
	int i := 0
	
	//slices := [][]int{}
	
	// checking 1-3 rows
	for i := 0; i < 3; i++ {
		for j := 0 + offset; j < 3; j++ {
			// validation logic
		}
	}
	
			// getting a valid index i
			validIndex := 
			// Putting the element at i into grid
			// Putting the element onto grid
			grid[i+offset][j+offset] = c[validIndex]
			// Put last element at index i
			c[randomIndex] = c[len(c)-1]
			// Decrement size of slice by 1
			c = c[:len(c)-1]
		}
	}
	wg.Done()
	return
}
*/
