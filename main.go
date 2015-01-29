package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	for {
		var m Metric
		if err := dec.Decode(&m); err != nil {
			log.Println(err)
			return
		}

		o := fmt.Sprintf("put %s %d %f", m.Name, m.Timestamp, m.Value)

		for name, value := range m.Tags {
			o += fmt.Sprintf(" %s=%s", name, value)
		}

		o += "\n"

		fmt.Print(o)
	}
}
