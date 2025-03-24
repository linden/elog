package elog

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
)

// Taken from https://github.com/linden/httplog/blob/a86d78bda99026f2520d02950b2fc4444a759ea5/httplog_test.go#L12-L65.
func parseLine(p []byte) map[string]string {
	// trim the line delimiter from the end.
	p = p[:len(p)-1]

	// create new map to store the components of the line.
	r := make(map[string]string)

	var k []byte
	var v []byte

	inV := false
	inS := false

	// iterate over remaining every charecter in the line.
	for _, c := range p {
		switch {
		// check if we've moved from the key to the value portion of the component.
		case c == '=' && inS == false:
			inV = true

		// check if we're starting a new key component.
		case c == ' ' && inS == false:
			inV = false

			// store the complete component in the map.
			r[string(k)] = string(v)

			// create new empty key and value.
			k = []byte{}
			v = []byte{}

		// check if we're starting or ending a string.
		case c == '"':
			if inS == true {
				inS = false
			} else {
				inS = true
			}

		// add to either the key of the value.
		default:
			if inV == true {
				v = append(v, c)
			} else {
				k = append(k, c)
			}
		}
	}

	// add the last key and value.
	r[string(k)] = string(v)

	return r
}

func TestElog(t *testing.T) {
	b := new(bytes.Buffer)

	// Create a structured logger.
	h := slog.NewTextHandler(b, nil)
	sl := slog.New(h)

	// Create a event logger.
	el := New(sl)

	// Create a new event.
	e := el.NewEvent("greeted")

	// Set the properties of the event.
	e.With("name", "Jim")
	e.WithError(errors.New("error"))
	e.WithWarn(errors.New("warn"))

	// Log the event to the buffer.
	e.Log()

	t.Logf("%s", b.String())

	// Parse the event.
	l := parseLine(b.Bytes())

	// Ensure the properties of the event match.
	if l["level"] != "ERROR" {
		t.Fatalf("expected level to be ERROR got %s", l["level"])
	}

	if l["msg"] != "greeted" {
		t.Fatalf("expected msg to be greeted got %s", l["msg"])
	}

	if l["name"] != "Jim" {
		t.Fatalf("expected name to be Jim got %s", l["name"])
	}

	if l["error"] != "error" {
		t.Fatalf("expected error to be error got %s", l["error"])
	}

	if l["warn"] != "warn" {
		t.Fatalf("expected warn to be warn got %s", l["warn"])
	}
}
