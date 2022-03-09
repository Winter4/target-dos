package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func launchThread(url string, index int, ch chan int) {

	for {
		time.Sleep(time.Millisecond * 1000)

		res, err := http.Get(url)

		if err != nil {
			log.Println(" Error on getting response in ", index, " thread")
			log.Fatal(err)
		} else {
			ch <- index
			ch <- res.StatusCode
		}
	}
}

func main() {

	url := getURL()
	clearTerminal()
	println("\tCurrent URL: ", url)

	ch := make(chan int)

	go launchThread(url, 0, ch)

	for i := range ch {
		log.Println(" |  Thread:", i, " Status:", <-ch)
	}
}
