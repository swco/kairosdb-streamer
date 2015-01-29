package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

// Metric holds the details of a metric
type Metric struct {
	Name      string            `json:"name"`
	Timestamp int               `json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

func main() {
	dec := json.NewDecoder(os.Stdin)

	conn, err := net.Dial("tcp", "localhost:4243")

	if err != nil {
		log.Fatalln(err)
	}

	for {
		var m Metric
		if err := dec.Decode(&m); err != nil {
			log.Fatalln(err)
		}

		o := fmt.Sprintf("put %s %d %f", m.Name, m.Timestamp, m.Value)

		for name, value := range m.Tags {
			o += fmt.Sprintf(" %s=%s", name, value)
		}

		o += "\r"

		fmt.Fprint(conn, o)
	}
}
