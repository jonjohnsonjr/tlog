package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if err := mainE(os.Stdout, os.Stdin); err != nil {
		log.Fatal(err)
	}
}

type line struct {
	when time.Time
	ts   string
	text string
}

func mainE(w io.Writer, r io.Reader) error {
	lines := []line{}
	scanner := bufio.NewScanner(r)

	start := time.Unix(1<<62, 0)
	end := time.Unix(0, 0)

	for scanner.Scan() {
		in := scanner.Text()
		before, after, ok := strings.Cut(in, " ")
		if !ok {
			log.Printf("no space in %q", in)
			continue
		}

		when, err := time.Parse(time.RFC3339, before)
		if err != nil {
			log.Printf("parsing timestamp %q: %v", before, err)
			continue
		}

		if when.Before(start) {
			start = when
		}
		if when.After(end) {
			end = when
		}

		lines = append(lines, line{
			when: when,
			ts:   before,
			text: after,
		})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading input: %w", err)
	}

	total := end.Sub(start)

	fmt.Fprint(w, header)

	for _, l := range lines {
		left := l.when.Sub(start)
		leftpad := float64(left) / float64(total)
		fmt.Fprintf(w, `<div>`)
		fmt.Fprintf(w, `<span style="background: linear-gradient(90deg, #EEEEEE %f%%, #FFFFFF %f%%);">%s %s</span>`, 100.0*leftpad, 100.0*leftpad, l.ts, l.text)
		fmt.Fprintf(w, "</div>\n")
	}

	fmt.Fprint(w, footer)
	return nil
}

// TODO
// background: linear-gradient(80deg, #ff0000 50%, #0000ff 50%);

const header = `
<html>
    <head>
        <title>tlog</title>
        <style>
        span {
          display: block;
          white-space: nowrap;
          font-family: monospace;
        }
        body {
        	display: block;
        }
        </style>
    </head>
    <body>`

const footer = `
    </body>
</html>
`
