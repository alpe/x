// logview formats JSON logs produced by x/log package and prints them to stdout.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	for {
		entry := make(map[string]string)
		if err := decode(&entry); err == io.EOF {
			return
		}

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
	var b bytes.Buffer

	fmt.Fprint(&b, FgGray)
	fmt.Fprint(&b, date.Format("15:04 "))
	fmt.Fprint(&b, NoColor)

	if len(msg) < 30 {
		msg += strings.Repeat(" ", 30-len(msg))
	}
	if level == "ERROR" {
		fmt.Fprint(&b, FgRed)
		fmt.Fprint(&b, msg)
		fmt.Fprint(&b, NoColor)
	} else {
		fmt.Fprint(&b, msg)
	}

	keys := make([]string, 0, len(extra))
	for k := range extra {
		if k == "file" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := extra[k]
		fmt.Fprintf(&b, " %s%s%s:%s ", FgGray, k, NoColor, v)
	}
	fmt.Fprintf(&b, " %sfile=%s%s", FgGray, extra["file"], NoColor)

	return b.String()
}

// console colors
const (
	NoColor = "\033[0m"

	FgBlack  = "\033[30m"
	FgRed    = "\033[31m"
	FgGreen  = "\033[32m"
	FgYellow = "\033[33m"
	FgBlue   = "\033[34m"
	FgGray   = "\033[38;05;239m"

	BgBlack  = "\033[40m"
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
	BgGray   = "\033[49m"
)
