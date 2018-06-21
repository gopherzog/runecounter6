package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const Version = "0.1.6"

type runeinfo map[string]int

type proverb struct {
	line  string
	chars map[string]int
}

func loadProverbs(path string) (error, []*proverb) {
	err, rawtext := readText(path)
	if err != nil {
		return err, nil
	}
	lines := strings.Split(rawtext, "\n")
	var proverbs []*proverb
	for _, line := range lines {
		runemap := mapRunes(line)
		myProverb := &proverb{line, runemap}
		proverbs = append(proverbs, myProverb)
	}
	return err, proverbs
}

func filenameFromArgs(args []string) (error, string) {
	filename := os.Getenv("FILE")
	maybePath := flag.String("f", filename, "file to be processed")
	flag.Parse()
	if *maybePath != "" {
		filename = *maybePath
	}
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return err, ""
		}
	}
	return nil, filename
}

func main() {
	fmt.Printf("runecounter2 %s\n", Version)
	err, filename := filenameFromArgs(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	err, proverbs := loadProverbs(filename)
	if err != nil {
		problem := fmt.Sprintf("%v", err)
		panic(problem)
	}

	ch0 := make(chan *proverb)
	ch1 := make(chan *proverb)
	go printProverbs(ch0, ch1)
	for idx, p := range proverbs {
		if idx%2 == 0 {
			ch0 <- p
		} else {
			ch1 <- p
		}
	}
	close(ch0)
	close(ch1)
}

func printProverbs(ch0, ch1 chan *proverb) {
	var p *proverb
	myPrinter := func(chanName string, p *proverb) {
		fmt.Printf("%s: %s\n", chanName, p.line)
		fmt.Printf("%s\n\n", formatMap(p.chars))
	}
	for {
		select {
		case p = <-ch0:
			myPrinter("First", p)

		case p = <-ch1:
			myPrinter("Second", p)
		}
	}
}

func mapRunes(line string) runeinfo {
	runes := make(map[string]int)
	for _, r := range line {
		s := string(r)
		runes[s] = runes[s] + 1
	}
	return runes
}

func formatMap(info runeinfo) string {
	var allItems []string
	for k, v := range info {
		item := fmt.Sprintf("'%s'=%d", k, v)
		allItems = append(allItems, item)
	}
	textRepr := strings.Join(allItems, ", ")
	return textRepr
}

func readText(filename string) (error, string) {
	dat, err := ioutil.ReadFile(filename)
	lx := len(dat)
	proverbs := ""
	if err == nil {
		maxChars := lx
		if string(dat[lx-1]) == "\n" {
			maxChars--
		}
		proverbs = string(dat[:maxChars])
	}
	return err, proverbs
}
