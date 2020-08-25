package main

import (	
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
)

// precalculated Seed values, Time Complexity will be reduced by a 
// large factor as each seed produces only a valid sudoku
// There are infinite valid seed values, however we will be dealing with limited of them
var seeds = []int{ 55,560,1890,4422,4585,5377,5583,6119,6886,7148,7295,8348,8847,8915,9518,9663,
		9848,11143,11194,11438,11709,11734,12495,13016,16300,16717,17027,17285,17541,17555}

var upgrader = websocket.Upgrader{}

type Sudoku struct {
	oneToNineArray	[9]int
	grid          	[9][9]int
	userGrid 	[9][9]int	
	level		int
}
// level 0 means easy
// level 1 means Medium
// level 2 means Hard


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
	s := Sudoku{}
	s.initializeSudoku()
	err = s.generateGrid()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.getGridForUser()
	
	fmt.Println("Answer: ")
	displayGrid(&s.grid)
	
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
		
		if checkForViolation(&s.userGrid, row, col) {
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
