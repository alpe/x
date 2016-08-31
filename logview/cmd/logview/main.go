// logview formats JSON logs produced by x/log package and prints them to stdout.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	entry := make(map[string]string)
	for scanner.Scan() {

		data := scanner.Bytes()
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			fmt.Println(string(data))
			continue
		}

		level := pop(entry, "level")
		msg := pop(entry, "msg")
		date, _ := time.Parse(time.RFC3339, pop(entry, "date"))
		file := pop(entry, "file")

		fmt.Fprintln(os.Stdout, fmtEntry(level, msg, date, file, entry))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

var decode = json.NewDecoder(os.Stdin).Decode

func pop(row map[string]string, key string) string {
	val := row[key]
	delete(row, key)
	return val
}

func fmtEntry(level, msg string, date time.Time, file string, extra map[string]string) string {
	chunks := make([]string, 0, 16)

	switch level {
	case "DEBUG":
		chunks = append(chunks, bgGreen("D"))
	case "ERROR":
		chunks = append(chunks, bgRed("E"))
	default:
		chunks = append(chunks, bgBlue("X"))
	}

	chunks = append(chunks, fgBlue(date.Format("15:04")))

	chunks = append(chunks, ":")

	if len(msg) < 30 {
		msg += strings.Repeat(" ", 30-len(msg))
	}
	chunks = append(chunks, msg)

	chunks = append(chunks, fgCyan(file))

	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := extra[k]
		chunks = append(chunks, fmt.Sprintf("%s%s%s=%s%s", FgBlue, k, FgBlue, NoColor, v))
	}

	return strings.Join(chunks, " ")
}

// console colors
const (
	NoColor = "\033[0m"

	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgGray    = "\033[38;05;237m"

	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgGray    = "\033[49m"
)

func fgBlack(params ...string) string {
	return FgBlack + strings.Join(params, " ") + NoColor
}
func fgRed(params ...string) string {
	return FgRed + strings.Join(params, " ") + NoColor
}

func fgGreen(params ...string) string {
	return FgGreen + strings.Join(params, " ") + NoColor
}
func fgYellow(params ...string) string {
	return FgYellow + strings.Join(params, " ") + NoColor
}
func fgBlue(params ...string) string {
	return FgBlue + strings.Join(params, " ") + NoColor
}

func fgMagenta(params ...string) string {
	return FgMagenta + strings.Join(params, " ") + NoColor
}

func fgCyan(params ...string) string {
	return FgCyan + strings.Join(params, " ") + NoColor
}

func bgBlack(params ...string) string {
	return BgBlack + strings.Join(params, " ") + NoColor
}

func bgRed(params ...string) string {
	return BgRed + strings.Join(params, " ") + NoColor
}

func bgGreen(params ...string) string {
	return BgGreen + strings.Join(params, " ") + NoColor
}

func bgYellow(params ...string) string {
	return BgYellow + strings.Join(params, " ") + NoColor
}
func bgBlue(params ...string) string {
	return BgBlue + strings.Join(params, " ") + NoColor
}
func bgMagenta(params ...string) string {
	return BgMagenta + strings.Join(params, " ") + NoColor
}
func bgGray(params ...string) string {
	return BgGray + strings.Join(params, " ") + NoColor
}
