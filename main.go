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

// Send writes the metric m into conn in the format expected by kairosdb's tcp api
func Send(w io.Writer, m Metric) {
	fmt.Fprintf(w, "put %s %d %f", m.Name, m.Timestamp, m.Value)

	for name, value := range m.Tags {
		//empty tags will generate an error on ingest
		if value != "" && name != "" {
			fmt.Fprintf(w, " %s=%s", name, value)
		}
	}

	fmt.Fprint(w, "\n")
}

// getInputReader opens a filename located at the argnum'th argument or so.Stdin if
// the argument does not exist
func getInputReader(argnum int) (io.ReadCloser, error) {
	if flag.NArg() > argnum {
		return os.Open(flag.Arg(argnum))
	}
	return os.Stdin, nil
}

func main() {
	host := flag.String("host", "localhost:4242", "The host:port to connect to. Defaults to 'localhost:4242'")
	ver := flag.Bool("version", false, "Print the version number and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("%v version: %s\n", os.Args[0], version)
		os.Exit(0)
	}

	input, err := getInputReader(0)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()

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

		Send(conn, m)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}
