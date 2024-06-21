package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uuid-disted/consumer/app"
)

func main() {
	brokersFile := flag.String("f", "brokers.txt", "The file containing host address of RabbitMQ brokers IP")
	flag.Parse()

	brokerHosts, err := readBrokersFile(*brokersFile)
	if err != nil {
		fmt.Printf("Error reading brokers file: %v\n", err)
		return
	}

	app := app.NewApplication(brokerHosts)
	go app.Run("uuids")

	app.StartServer(8080)
}

func readBrokersFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	brokers := strings.Split(strings.TrimSpace(string(data)), "\n")
	return brokers, nil
}
