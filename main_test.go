package main

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createListener(t *testing.T, done chan int, port, expected string) net.Listener {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		t.Fatal("Could not setup listener:", err)
	}

	// Only expect one connection, then can shutdown the listener
	go func() {
		defer ln.Close()
		c, err := ln.Accept()
		if err != nil {
			t.Fatal("Error connecting to client:", err)
		}
		defer c.Close()

		buf, err := ioutil.ReadAll(c)
		if err != nil {
			t.Fatal("Error reading from client:", err)
		}

		assert.Equal(t, expected, string(buf))

		done <- 1
	}()

	return ln
}

func TestMetricGetsWrittenToConn(t *testing.T) {
	expected := "put test1 1429001359824 5.000000 tag1=bar tag2=foo\n"
	m := Metric{"test1", 1429001359824, 5, map[string]string{"tag1": "bar", "tag2": "foo"}}
	done := make(chan int)
	createListener(t, done, ":54000", expected)
	conn, err := net.Dial("tcp", ":54000")
	if err != nil {
		t.Fatal("Could not connect to listener: ", err)
	}

	Send(conn, m)
	conn.Close()
	<-done
}
