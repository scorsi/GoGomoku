package com

import (
	"gogomoku/ai"
	"os"
	"bufio"
	"sync"
	"fmt"
	"time"
	"regexp"
	"strconv"
)

type ComFunc func(string)

type ComFuncTab struct {
	fun		ComFunc
	reg		string
}

var comFuncTab = [6]ComFuncTab{
	{ fun: restartGame, reg: "RESTART"},
	{ fun: launchAI, reg: "BEGIN" },
	{ fun: startGame, reg: "^START" },
	{ fun: enemyTurn, reg: "TURN" },
	{ fun: endGame, reg: "END" },
	{fun: aboutAI, reg: "ABOUT"},
}

var isActive bool = true

func aboutAI(com string) {
	fmt.Println("name=GoGomoku, version=1.0, author=SaltTeam, country=France")
}

func endGame(_ string) {
	isActive = false
}

func restartGame(com string) {
	for line := range ai.GameBoard {
		for section := range  ai.GameBoard[line] {
			ai.GameBoard[line][section] = 0
		}
	}
	fmt.Println("OK - everything is good")
}

func enemyTurn(com string)  {
	r, err := regexp.Compile("TURN ([0-9]+),([0-9]+)")
	if err != nil {
		fmt.Println(err)
	}
	s := r.FindStringSubmatch(com)
	x, err := strconv.Atoi(s[1])
	if err != nil {
		fmt.Println(err)
	}
	y, err := strconv.Atoi(s[2])
	if err != nil {
		fmt.Println(err)
	}

	if len(ai.GameBoard) >= x && len(ai.GameBoard[x - 1]) >= y {
		ai.GameBoard[x - 1][y - 1] = 1
	} else {
		fmt.Println("Index out of range")
	}
	fmt.Println("OK - everything is good")
	ai.LastEnemyPosition = ai.Position{x, y}
	launchAI(com)
}

func startGame(com string) {
	r, err := regexp.Compile("START ([0-9]+)")
	if err != nil {
		fmt.Println(err)
	}
	size, err := strconv.Atoi(r.FindStringSubmatch(com)[1])
	if err != nil {
		fmt.Println(err)
	}
	ai.GameBoard = make([][]int, size)
	for x := range ai.GameBoard {
		ai.GameBoard[x] = make([]int, size)
		for y := range ai.GameBoard[x] {
			ai.GameBoard[x][y] = 0
		}
	}
	fmt.Println("OK - everything is good")
}

func launchAI(_ string) {
	wg := new(sync.WaitGroup)
	comChan := make(chan string, 1)

	wg.Add(2)
	go func() {
		defer wg.Done()

		ai.Start(comChan)
	}()
	time.Sleep(time.Millisecond * 5)
	go func() {
		defer wg.Done()

		select {
		case res := <- comChan:
			fmt.Println(res)
		case <- time.After(time.Second * 1):
			ai.StartRandom(comChan)
			fmt.Println(<- comChan)
		}
	}()
	wg.Wait()
	close(comChan)
}

func parseCom(com string) {
	for _, elem := range comFuncTab {
		match, err := regexp.Match(elem.reg, []byte(com))
		if err != nil {
			fmt.Print(err)
		} else if match {
			elem.fun(com)
		}
	}
}

func ComManagement() {
	com := new(ComStruct)
	com.reader = bufio.NewReader(os.Stdin)

	isActive = true

	for isActive {
		msg, err := com.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		parseCom(msg)
	}
}
