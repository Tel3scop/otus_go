package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	if err := Copy(from, to, offset, limit); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printProgressBar(current, total int64) {
	const barWidth = 50
	percent := float64(current) / float64(total)
	bar := int(percent * barWidth)
	fmt.Printf("\r[%s%s] %3.2f%%", strings.Repeat("=", bar), strings.Repeat(" ", barWidth-bar), percent*100)
}
