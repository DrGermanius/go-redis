package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	RESPseparator    = "\r\n"
	RESPseparatorLen = 2
)

func processRESPArray(b []byte) []byte {
	elems := strings.Split(string(b[:len(b)-RESPseparatorLen]), RESPseparator)

	switch elems[2] {
	case "ECHO", "echo":
		var res string
		for i := 4; i < len(elems); i = i + 2 {
			res += elems[i] + " "
		}
		return encodeToRESPSimpleString(res[:len(res)-1])
	case "PING", "ping":
		return encodeToRESPSimpleString("PONG")
	case "SET", "set":
		i := Item{
			Value: elems[6],
		}
		if len(elems) > 8 && (elems[8] == "PX" || elems[8] == "px") {
			exp, _ := strconv.Atoi(elems[10])
			i.Deadline = time.Now().Add(time.Millisecond * time.Duration(exp))
		}
		cache.Add(elems[4], i)
		return encodeToRESPSimpleString("OK")
	case "GET", "get":
		el, ok := cache.Get(elems[4])
		if !ok {
			return nullString()
		}
		if time.Now().After(el.Deadline) && !el.Deadline.IsZero() {
			cache.Remove(elems[4])
			return nullString()
		}
		return encodeToRESPSimpleString(el.Value)
	}
	return nullString()
}

func isRESPArray(b []byte) bool {
	return b[0] == 42 // RESP array starts with '*'
}

func encodeToRESPSimpleString(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}

func nullString() []byte {
	return []byte("$-1\r\n")
}
