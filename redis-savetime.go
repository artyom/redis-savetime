package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(0)
	var (
		start time.Time
		end   time.Time
		err   error
	)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s /path/to/redis.log", os.Args[0])
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		switch {
		case bytes.Contains(scanner.Bytes(), savingStart):
			start, err = readTime(scanner.Bytes())
		case bytes.Contains(scanner.Bytes(), savingEnd):
			end, err = readTime(scanner.Bytes())
		default:
			continue
		}
		if err != nil {
			log.Print(err)
			continue
		}
		if start.IsZero() || end.IsZero() || end.Before(start) {
			continue
		}
		fmt.Printf("%s\t%s\t%s\n", end.Sub(start),
			start.Format(time.StampMilli),
			end.Format(time.StampMilli),
		)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var (
	savingStart = []byte(`Background saving started`)
	savingEnd   = []byte(`Background saving terminated`)
)

const layout = `_2 Jan 15:04:05.000`

func readTime(line []byte) (time.Time, error) {
	b := bytes.Index(line, []byte(`] `))
	if b < 0 {
		return time.Time{}, fmt.Errorf("cannot find date beginnig")
	}
	e := bytes.Index(line, []byte(` *`))
	if e < 0 {
		return time.Time{}, fmt.Errorf("cannot find date end")
	}
	return time.Parse(layout, string(line[b+2:e]))
}
