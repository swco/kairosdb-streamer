package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleJSON = `{"timestamp":1429001359824,"name":"test.metric1","value":42.0,"tags":{"tag1":"value1","tag2":"value2"}}
{"timestamp":1429001359824,"name":"test.metric2","value":43.0,"tags":{"tag1":"value1"}}
{"timestamp":1429001359824,"name":"test.metric3","value":42.0}`

func TestDecodeJSONIntoMetric(t *testing.T) {

}

func TestInvalidJSONGetsDropped(t *testing.T) {

}

func TestMetricGetsWrittenToConn(t *testing.T) {
	metrics := []Metric{
		Metric{"test1", 1429001359824, 5, map[string]string{"tag1": "bar", "tag2": "foo"}},
		Metric{"test2", 1429001359824, 54, map[string]string{}},
		Metric{"test3", 1429001359824, 1, map[string]string{"tag1": "bar"}},
		Metric{"test4", 1429001359824, 5, map[string]string{"tag1": "bar", "tag2": ""}},
		Metric{"test5", 1429001359824, 5, map[string]string{"tag1": "baz", "": "bar", "tag2": ""}},
	}

	expected := []Metric{
		Metric{"test1", 1429001359824, 5, map[string]string{"tag1": "bar", "tag2": "foo"}},
		Metric{"test2", 1429001359824, 54, map[string]string{}},
		Metric{"test3", 1429001359824, 1, map[string]string{"tag1": "bar"}},
		Metric{"test4", 1429001359824, 5, map[string]string{"tag1": "bar"}},
		Metric{"test5", 1429001359824, 5, map[string]string{"tag1": "baz"}},
	}

	rw := &bytes.Buffer{}

	for _, m := range metrics {
		Send(rw, m)
	}

	s := bufio.NewScanner(rw)
	for _, m := range expected {
		if !s.Scan() {
			assert.Fail(t, "Output has one line for each valid metric:", s.Err())
		}

		line := s.Text()
		parts := strings.Split(line, " ")

		assert.True(t, len(parts) >= 4, "Output has more then 4 fields")
		assert.Equal(t, "put", parts[0], "Output starts with 'puts'")
		assert.Equal(t, m.Name, parts[1], "Outputs second field is the metrics name")
		assert.Equal(t, fmt.Sprintf("%d", m.Timestamp), parts[2], "Outputs third field is the metrics timestamp")
		assert.Equal(t, fmt.Sprintf("%f", m.Value), parts[3], "Outputs fourth field is the metrics value")

		tags := map[string]string{}
		for _, part := range parts[4:] {
			tag := strings.Split(part, "=")
			assert.True(t, len(tag) == 2, "Tag has a key and value,", part)
			if len(tag) == 2 {
				tags[tag[0]] = tag[1]
			}
		}

		assert.Equal(t, m.Tags, tags, "Output contains all tags")
	}

	if err := s.Err(); err != nil {
		t.Fatal("Error reading from client:", err)
	}
}
