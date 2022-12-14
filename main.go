package main

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type NodeKasami struct {
	RN       []int
	hasToken bool
}

type SuzukiKasami struct {
	nodes    []NodeKasami
	LN       []int
	Q        *list.List
	processN int
}

func releaseSC(processID int, sk *SuzukiKasami) {
	sk.LN[processID] = sk.nodes[processID].RN[processID]

	for i := 0; i < sk.processN; i++ {
		if (!find(i, sk.Q)) && (sk.nodes[processID].RN[i] == sk.LN[i]+1) {
			sk.Q.PushBack(i)
		}
	}

	if sk.Q.Front() != nil {
		sk.nodes[processID].hasToken = false
		sk.nodes[sk.Q.Front().Value.(int)].hasToken = true
		sk.Q.Remove(sk.Q.Front())
	}
}

func requestSC(processID int, sk *SuzukiKasami) {
	sk.nodes[processID].RN[processID]++

	// Solicitar a los demás procesos
	for i := 0; i < sk.processN; i++ {
		if i == processID {
			continue
		}

		sk.nodes[i].RN[processID] = int(math.Max(float64(sk.nodes[processID].RN[processID]), float64(sk.nodes[i].RN[processID])))
		if sk.nodes[i].RN[processID] == sk.LN[processID]+1 && (sk.nodes[i].hasToken && (sk.LN[i] == sk.nodes[i].RN[i])) {
			sk.nodes[i].hasToken = false
			sk.nodes[processID].hasToken = true
		}
	}
}

// Auxiliar
func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func find(i int, Q *list.List) bool {
	for e := Q.Front(); e != nil; e = e.Next() {
		if i == e.Value {
			return true
		}
	}
	return false
}

func compare(RN []int, LN []int, size int) bool {
	for i := 0; i < size; i++ {
		if RN[i] == (LN[i] + 1) {
			return false
		}
	}
	return true
}

func getPassanger() (string, string, bool) {
	file, err := os.ReadFile("pasajeros.txt")

	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(file), "\n")

	if lines[0] == "" || len(lines) == 0 {
		return "EOF", "EOF", true
	}

	passenger := strings.Split(lines[0], " ")

	lastName := passenger[0]
	seat := passenger[1]

	lines = remove(lines, 0)

	output := strings.Join(lines, "\n")
	err = os.WriteFile("pasajeros.txt", []byte(output), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	return lastName, seat, false
}

func updateLog(processID int, lastName string, seat string) {
	logFile, err := os.OpenFile("procesados.txt", os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer logFile.Close()

	_, err = fmt.Fprintln(logFile, "P"+strconv.Itoa(processID+1)+": "+lastName+" "+seat)

	if err != nil {
		fmt.Println(err)
	}
}

func updateMap(number string) {
	file, err := os.ReadFile("mapa.txt")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(file, []byte(strings.TrimSpace(number)), []byte("XX"), 1)

	if err = os.WriteFile("mapa.txt", output, 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func updateProfits(seatString string) {
	file, err := os.Open("ganancias.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	profits := make([]int, 3)

	for index := 0; scanner.Scan(); index++ {
		profitText := scanner.Text()
		profit, _ := strconv.Atoi(profitText)

		profits[index] = profit
	}

	seat, err := strconv.Atoi(strings.TrimSpace(seatString))

	if seat >= 1 && seat <= 16 {
		profits[0] += 8000
	} else if seat >= 17 && seat <= 32 {
		profits[1] += 6000
	} else if seat >= 33 && seat <= 48 {
		profits[2] += 4000
	}

	output := strings.Trim(fmt.Sprint(profits), "[]")

	if err = os.WriteFile("ganancias.txt", []byte(output), 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processPassanger(wg *sync.WaitGroup, processID int, sk *SuzukiKasami) {
	requested := false
	for {
		/*
			El TOKEN siempre lo tendrá el primer proceso, por lo que
			los demás si lo necesitan deberán pedirlo desde este punto

			EJ (3 procesos):
				R_1 = {0, 0, 0}
				L = {0, 0, 0}
				Q = {}
		*/
		// ejecutar SC
		if sk.nodes[processID].hasToken {
			// Lectura y modificación de archivos
			lastName, seat, flag := getPassanger()

			if flag {
				releaseSC(processID, sk)
				break
			}

			updateMap(seat)
			updateProfits(seat)
			updateLog(processID, lastName, seat)

			time.Sleep(time.Second)
			releaseSC(processID, sk)
			requested = false
		} else if !requested {
			requestSC(processID, sk)
			requested = true
		}
	}

	wg.Done() // Resta uno al contador
}

func main() {
	// limpiar log
	logFile, err := os.Create("procesados.txt")

	if err != nil {
		log.Fatal(err)
	}

	logFile.Close()

	var sk SuzukiKasami

	sk.processN, _ = strconv.Atoi(os.Args[1])
	sk.nodes = make([]NodeKasami, sk.processN)

	for i := 0; i < sk.processN; i++ {
		sk.nodes[i].RN = make([]int, sk.processN)
		for j := 0; j < sk.processN; j++ {
			sk.nodes[i].RN[j] = 0
		}
	}

	sk.LN = make([]int, sk.processN)
	sk.Q = list.New()
	var wg sync.WaitGroup

	// Always first process has token
	sk.nodes[0].hasToken = true

	// Logica de procesos
	for i := 0; i < sk.processN; i++ {
		wg.Add(1) // Añade uno al contador
		go processPassanger(&wg, i, &sk)
	}

	wg.Wait()
}
