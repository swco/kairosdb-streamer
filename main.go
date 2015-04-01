package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

const version = "0.0.2"

// Metric holds the details of a metric
type Metric struct {
	Name      string            `json:"name"`
	Timestamp int               `json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

func main() {
	host := flag.String("host", "localhost:4242", "The host:port to connect to. Defaults to 'localhost:4242'")
	ver := flag.Bool("version", false, "Print the version number and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("%v version: %s\n", os.Args[0], version)
		os.Exit(0)
	}

	inputFilename := flag.Arg(0)

	var input io.Reader

	if len(inputFilename) > 0 {
		var err error

		input, err = os.Open(inputFilename)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		input = os.Stdin
	}

	dec := json.NewDecoder(input)
	conn, err := net.Dial("tcp", *host)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		var m Metric
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		o := fmt.Sprintf("put %s %d %f", m.Name, m.Timestamp, m.Value)

		for name, value := range m.Tags {
			o += fmt.Sprintf(" %s=%s", name, value)
		}

		o += "\n"

		fmt.Fprint(conn, o)
	}
}
