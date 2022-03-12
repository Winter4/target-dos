package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getURL() string {

	var input string
	reader := bufio.NewReader(os.Stdin)

	for {
		print("\nType URL: ")

		input, _ = reader.ReadString('\n')
		input = strings.Replace(input, "\r\n", "", -1)
		if input == "" {
			continue
		}

		println("\nYou typed: ", input)
		print("Is it right? \n Type Y/N > ")

		confirm, _ := reader.ReadString('\n')
		confirm = strings.Replace(confirm, "\r\n", "", -1)

		if confirm == "Y" || confirm == "y" {
			break
		}
	}

	return "http://" + input
}

func launchThread(url string, index int, ch chan [2]string) {

	arr := [2]string{}
	for {
		res, err := http.Get(url)

		arr[0] = strconv.Itoa(index)
		if err != nil {
			log.Println(" Error on getting response in ", index, " thread")
			log.Println(err)

			arr[1] = strconv.Itoa(-1)
		} else {
			arr[1] = strconv.Itoa(res.StatusCode)
		}

		ch <- arr
	}
}

func printMatrix(info [][]int) {

	for i := range info {
		log.Println("  |  Thread: ", i+1, "  Status: ", info[i][1])
	}
}

func handleResponse(ch chan [2]string, info [][]int) {

	for range ch {
		arr := <-ch
		i, _ := strconv.Atoi(arr[0])
		code, _ := strconv.Atoi(arr[1])

		info[i][1] = code
	}
}

func main() {

	url := getURL()
	clearTerminal()

	ch := make(chan [2]string)

	// info matrix allocating
	threadsNum := runtime.NumCPU() * 5
	info := make([][]int, threadsNum)
	for i := 0; i < threadsNum; i++ {
		info[i] = make([]int, 2)
		info[i][0] = i

		go launchThread(url, i, ch)
	}

	// ________________________________________

	// main loop
	go handleResponse(ch, info)
	for {
		clearTerminal()
		println("\tCurrent URL: ", url)
		printMatrix(info)

		time.Sleep(time.Second * 2)
	}
}

// TODO: error handling
