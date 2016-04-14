// logview formats JSON logs produced by x/log package and prints them to stdout.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	for {
		entry := make(map[string]string)
		decode(&entry)

		level := pop(entry, "level")
		msg := pop(entry, "msg")
		date, _ := time.Parse(time.RFC3339, pop(entry, "date"))
		fmt.Fprintln(os.Stdout, fmtEntry(level, msg, date, entry))
	}
}

var decode = json.NewDecoder(os.Stdin).Decode

func pop(row map[string]string, key string) string {
	val := row[key]
	delete(row, key)
	return val
}

func fmtEntry(level, msg string, date time.Time, extra map[string]string) string {
	chunks := make([]string, 0, 16)

	switch level {
	case "DEBUG":
		chunks = append(chunks, BgGreen, "D", NoColor)
	case "ERROR":
		chunks = append(chunks, BgRed, "E", NoColor)
	default:
		chunks = append(chunks, BgBlue, "X", NoColor)
	}

	chunks = append(chunks, FgGray, date.Format("15:04"), NoColor)

	if len(msg) < 30 {
		msg += strings.Repeat(" ", 30-len(msg))
	}
	chunks = append(chunks, msg)

	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := extra[k]
		chunks = append(chunks, fmt.Sprintf("%s%s%s=%s%s", FgGray, k, FgBlack, NoColor, v))
	}

	return strings.Join(chunks, " ")
}

// console colors
const (
	NoColor = "\033[0m"

	FgBlack  = "\033[30m"
	FgRed    = "\033[31m"
	FgGreen  = "\033[32m"
	FgYellow = "\033[33m"
	FgBlue   = "\033[34m"
	FgGray   = "\033[38;05;237m"

	BgBlack  = "\033[40m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgGray   = "\033[49m"
)
