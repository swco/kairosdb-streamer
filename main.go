package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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

// send writes the metric m into conn in the format expected by kairosdb's tcp api
func send(w io.Writer, m Metric) {
	fmt.Fprintf(w, "put %s %d %f", m.Name, m.Timestamp, m.Value)

	for name, value := range m.Tags {
		//empty tags will generate an error on ingest
		if value != "" && name != "" {
			fmt.Fprintf(w, " %s=%s", name, value)
		}
	}

	fmt.Fprint(w, "\n")
}

func valid(m Metric) bool {
	return len(m.Name) != 0 && m.Timestamp != 0
}

// process reads metrics in json format from in and writes them to out in the
// format needed by kairosdb's tcp api
func process(in io.Reader, out io.Writer, errorLog *log.Logger) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		var m Metric

		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			errorLog.Printf("unable to decode line: %s: '%s'\n", err.Error(), scanner.Text())
			continue
		}

		if valid(m) {
			send(out, m)
		}
	}
	return scanner.Err()
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
	var errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	host := flag.String("host", "localhost:4242", "The host:port to connect to. Defaults to 'localhost:4242'")
	ver := flag.Bool("version", false, "Print the version number and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("%v version: %s\n", os.Args[0], version)
		os.Exit(0)
	}

	in, err := getInputReader(0)
	if err != nil {
		errorLog.Println(err)
		os.Exit(1)
	}
	defer in.Close()

	out, err := net.Dial("tcp", *host)
	if err != nil {
		errorLog.Println(err)
		os.Exit(1)
	}
	defer out.Close()

	if err := process(in, out, errorLog); err != nil {
		errorLog.Println("error reading input:", err)
		os.Exit(1)
	}
}
