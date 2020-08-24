package main

import (	
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

// precalculated Seed values, Time Complexity will be reduced by a large 
// factor as each seed only produces a valid sudoku
var seeds = []int{ 55,560,1890,4422,4585,5377,5583,6119,6886,7148,7295,8348,8847,8915,9518,9663,
		9848,11143,11194,11438,11709,11734,12495,13016,16300,16717,17027,17285,17541,17555}

var upgrader = websocket.Upgrader{}

// Blank Boxes for userGrid
var blankBoxes int = 5

type Sudoku struct {
	oneToNineArray	[9]int
	grid          	[9][9]int
	userGrid 	[9][9]int	
}



func main() {
	router := InitRouter()
	server := negroni.Classic()
	server.UseHandler(router)
	
	server.Run(":12345")
}

func InitRouter() (router *mux.Router){
	router = mux.NewRouter()
	
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./assets/"))))
	router.HandleFunc("/", homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/ws", newGameHandler).Methods(http.MethodGet)
	
	return
}

func homeHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "game.html")
}

// Request Handler for new Game
func newGameHandler(rw http.ResponseWriter, req *http.Request){
	c, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		log.Print("Upgrade: ", err)
	}
	// c.WriteMessage(websocket.TextMessage, []byte("Hello from Server"))
	s := Sudoku{}
	s.initializeSudoku()
	err = s.generateGrid()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.getGridForUser(blankBoxes)
	
	fmt.Println("Answer: ")
	displayGrid(s.grid)
	
	str := getStringArray( s.userGrid )
	fmt.Println("Server String: ", str)
	
	c.WriteMessage(websocket.TextMessage , []byte(str))
	
	for {
		_, recvData, err := c.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		
		//Extracting data from UI
		data := string(recvData)
		split := strings.Split(data, ",")
		value, _ := strconv.Atoi(split[0])
		row, _ := strconv.Atoi(split[1])
		col, _ := strconv.Atoi(split[2])
		
		s.userGrid[row][col] = value
		if s.grid[row][col] == value {
			win := s.checkWin()
			if win{
				c.WriteMessage(websocket.TextMessage, []byte("win"))
				break
			}			
		} else {
			c.WriteMessage(websocket.TextMessage, []byte("violation"))
		}	
	}	
}


func (s *Sudoku) initializeSudoku() {

	// Initialization Logic
	
	for i := 0; i < 9; i++ {
		s.oneToNineArray[i] = i + 1
		// initializing the matrix to blank := 0
		for j := 0; j < 9; j++ {
			s.grid[i][j] = 0
			s.userGrid[i][j] = 0
		}
	}
}
	

func(s *Sudoku) generateGrid() error {
	//
	var wg sync.WaitGroup

	// seedValue is what we pass to random function for first three block generation
	r := rand.New(rand.NewSource(time.Now().UnixNano() ) ) 
	randomValue :=r.Intn( len(seeds) - 1 )	// 29 is the length of Valid seedValues
	

	seedValue := seeds[randomValue]

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
		go s.parallelBlock048(offset, seedValue, &wg)
	}

	// wait for the threads
	wg.Wait()

	//
	// Logic which fills block 2, 3, 7 of the suduko
	//
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

	// Logic for block 1, 5, 6

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
	value := s.CompleteValidation()
	if value == true {
		//fmt.Print(seedValue)
		return nil
	} else {
		return fmt.Errorf("Sodoku wasn't generated!")
	}
	// Validation logic
	//sudukoValidation := CompleteValidation()

	//fmt.Println(sudukoValidation)
}


func (s *Sudoku) getGridForUser(blankBoxes int){
	s.userGrid = s.grid
	r := rand.New(rand.NewSource(time.Now().UnixNano() ) ) 
	for i:=0 ; i < blankBoxes; i++ {
		row := r.Intn(9)
		col := r.Intn(9)
		s.userGrid[row][col] = 0	
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


func (s *Sudoku) checkWin() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.grid[i][j] != s.userGrid[i][j] {
				return false
			}
		}
	}
	return true
}


// displayGrid()


func displayGrid(grid [9][9]int) {

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
	// Printing the grid
	
	fmt.Println()
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print("\t", grid[i][j])
		}
		fmt.Println()
	}
}


func  (s *Sudoku) parallelBlock048(offset, seedValue int, wg *sync.WaitGroup) {

	b := s.oneToNineArray // b is copy of array
	c := b[:]           // c is a slice referring array b
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

// needs editing
func  (s *Sudoku) parallelBlock237(blockID int, wg *sync.WaitGroup) {

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
			s.enterRow(posX+i, mapIndex, mapp)
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
		}
	}

	// inserting values into Suduko

	//value :=
	s.insert(blockID, mapp)
	/*if value == false {
		//fmt.Println("Suduko not possible, Try with other seed")
		os.Exit(-1)
	}*/

	wg.Done()
	return
}

func  (s *Sudoku) insert(blockID int, mapp map[int][]int) bool {

	arrayOfSlices := [9][]int{}
	sequence := [9]int{}
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

	// Adding elements

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

func isSafe(sequence *[9]int, arrayOfSlices [9][]int, k, index, blockID int) bool {
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

	if len(mapp) == k+1 { // && dependentBlock( blockID, sequence ){
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

// Validation Logic
func(s *Sudoku) CompleteValidation() bool {

	//truthValue := true

	truthValue := s.validateRows() && s.validateColumns() && s.validateBlocks()

	return truthValue

}

func(s *Sudoku) validateRows() bool {
	for i := 0; i < 9; i++ {
		if truthValueRow := s.validateSingleRow(i); truthValueRow == false {
			return false
		}
	}
	return true
}

func(s *Sudoku) validateSingleRow(RowNum int) bool {
	validationMap := make(map[int]bool, 9)
	for j := 0; j < 9; j++ {
		if s.grid[RowNum][j] >= 1 && s.grid[RowNum][j] <= 9 {
			validationMap[s.grid[RowNum][j]] = true
		} else {
			return false
		}
	}

	if len(validationMap) == 9 {
		return true
	}
	return false
}

func(s *Sudoku) validateColumns() bool {
	for i := 0; i < 9; i++ {
		if truthValueColumn := s.validateSingleColumn(i); truthValueColumn == false {
			return false
		}
	}
	return true
}

func(s *Sudoku) validateSingleColumn(ColumnNum int) bool {

	validationMap := make(map[int]bool, 9)
	for j := 0; j < 9; j++ {
		if s.grid[j][ColumnNum] >= 1 && s.grid[j][ColumnNum] <= 9 {
			validationMap[s.grid[j][ColumnNum]] = true
		} else {
			return false
		}
	}

	if len(validationMap) == 9 {
		return true
	}
	return false
}

func(s *Sudoku) validateBlocks() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Multilying by 3 , so as to get proper offset
			if truthValueBlock := s.validateSingleBlock(i*3, j*3); truthValueBlock == false {
				return false
			}
		}
	}
	return true
}

func(s *Sudoku) validateSingleBlock(offsetX, offsetY int) bool {

	validationMap := make(map[int]bool, 9)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if s.grid[offsetX+i][offsetY+j] >= 1 && s.grid[offsetX+i][offsetY+j] <= 9 {
				validationMap[s.grid[offsetX+i][offsetY+j]] = true
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
