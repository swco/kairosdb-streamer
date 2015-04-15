package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

const version = "0.0.3"

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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	} else {
		input = os.Stdin
	}

	scanner := bufio.NewScanner(input)
	conn, err := net.Dial("tcp", *host)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	for scanner.Scan() {
		var m Metric

		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			fmt.Fprintf(os.Stderr, "unable to decode line: %s: '%s'\n", err.Error(), scanner.Text())
			continue
		}

		o := fmt.Sprintf("put %s %d %f", m.Name, m.Timestamp, m.Value)

		for name, value := range m.Tags {
			//empty tags will generate an error on ingest
			if value != "" {
				o += fmt.Sprintf(" %s=%s", name, value)
			}
		}

		o += "\n"

		fmt.Fprint(conn, o)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}
