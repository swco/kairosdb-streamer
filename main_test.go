package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runTest(t *testing.T, inJSON string, expected []Metric) {
	in := strings.NewReader(inJSON)
	out := &bytes.Buffer{}

	err := process(in, out, log.New(ioutil.Discard, "", 0))
	assert.NoError(t, err, "Process should not produce an error")

	s := bufio.NewScanner(out)
	for _, m := range expected {
		if !s.Scan() {
			assert.Fail(t, "Output has one line for each valid metric")
			break
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

	assert.False(t, s.Scan(), "There are no more lines in the input")

	if err := s.Err(); err != nil {
		t.Fatal("Error reading from client:", err)
	}
}

func TestValidMetricsAppearInOutput(t *testing.T) {
	runTest(t,
		`{"timestamp":1429001359824,"name":"test1","value":1.0,"tags":{"tag1":"bar","tag2":"foo"}}
{"timestamp":1429001359824,"name":"test2","value":2.0,"tags":{"tag1":"bar","tag2":"foo"}}
{"timestamp":1429001359824,"name":"test3","value":2.0}
{"timestamp":1429001359824,"name":"test4","value":2}
{"timestamp":1429001359824,"name":"test5","value":22}
{"timestamp":1429001359824,"name":"test6","value":-43}
{"timestamp":1429001359824,"name":"test7","value":-44.123456789}
{"timestamp":1429001359824,"name":"test8","value":9876543210}
{"timestamp":1429001359824,"name":"test9","value":3.0001}`,
		[]Metric{
			Metric{"test1", 1429001359824, 1, map[string]string{"tag1": "bar", "tag2": "foo"}},
			Metric{"test2", 1429001359824, 2, map[string]string{"tag1": "bar", "tag2": "foo"}},
			Metric{"test3", 1429001359824, 2, map[string]string{}},
			Metric{"test4", 1429001359824, 2, map[string]string{}},
			Metric{"test5", 1429001359824, 22, map[string]string{}},
			Metric{"test6", 1429001359824, -43, map[string]string{}},
			Metric{"test7", 1429001359824, -44.123456789, map[string]string{}},
			Metric{"test8", 1429001359824, 9876543210, map[string]string{}},
			Metric{"test9", 1429001359824, 3.0001, map[string]string{}},
		})
}

func TestEmptyTagValueGetsDropped(t *testing.T) {
	runTest(t,
		`{"timestamp":1429001359824,"name":"test1","value":1.0,"tags":{"tag1":"bar","tag2":""}}
{"timestamp":1429001359824,"name":"test2","value":2.0,"tags":{"tag1":"","tag2":"foo"}}
{"timestamp":1429001359824,"name":"test3","value":3.0,"tags":{"tag1":"","tag2":""}}`,
		[]Metric{
			Metric{"test1", 1429001359824, 1, map[string]string{"tag1": "bar"}},
			Metric{"test2", 1429001359824, 2, map[string]string{"tag2": "foo"}},
			Metric{"test3", 1429001359824, 3, map[string]string{}},
		})
}

func TestEmptyTagKeyGetsDropped(t *testing.T) {
	runTest(t,
		`{"timestamp":1429001359824,"name":"test1","value":1.0,"tags":{"tag1":"bar","":"foo"}}
{"timestamp":1429001359824,"name":"test2","value":2.0,"tags":{"":"bar","tag2":"foo"}}
{"timestamp":1429001359824,"name":"test3","value":3.0,"tags":{"":"bar","":"foo"}}`,
		[]Metric{
			Metric{"test1", 1429001359824, 1, map[string]string{"tag1": "bar"}},
			Metric{"test2", 1429001359824, 2, map[string]string{"tag2": "foo"}},
			Metric{"test3", 1429001359824, 3, map[string]string{}},
		})
}

func TestInvalidJSONGetsDropped(t *testing.T) {
	runTest(t,
		`{"timestamp":1429001359824,"name":"test1","value":1.0,"tags":{"tag1":"bar"}}
{"timestamp":1429001359824,"name":"test2","value":1.0,"tags":{"tag1":"bar"}
"timestamp":1429001359824,"name":"test3","value":1.0,"tags":{"tag1":"bar"}}
{"timestamp"1429001359824,"name":"test4","value":1.0,"tags":{"tag1":"bar"}}
{"timestamp":1429001359824,name":"test5","value":1.0,"tags":{"tag1":"bar"}}

{"timestamp":1429001359824,"name":"test7""value":1.0,"tags":{"tag1":"bar"}}
{"timestamp":1429001359824,"name":"test8","value":1.0,"tags":{"tag1":"bar"}}`,
		[]Metric{
			Metric{"test1", 1429001359824, 1, map[string]string{"tag1": "bar"}},
			Metric{"test8", 1429001359824, 1, map[string]string{"tag1": "bar"}},
		})
}

func TestMissingRequiredKeysGetDropped(t *testing.T) {
	runTest(t,
		`{"timestamp":1429001359824,"name":"test1","value":1.0,"tags":{"tag1":"bar"}}
{"name":"test2","value":1.0,"tags":{"tag1":"bar"}}
{"timestamp":1429001359824,"value":1.0,"tags":{"tag1":"bar"}}
{}
{"timestamp":1429001359824,"name":"test2","value":1.0,"tags":{"tag1":"bar"}}`,
		[]Metric{
			Metric{"test1", 1429001359824, 1, map[string]string{"tag1": "bar"}},
			Metric{"test2", 1429001359824, 1, map[string]string{"tag1": "bar"}},
		})
}
