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

	ch := make(chan *proverb)
	go printProverbs(ch)
	for _, p := range proverbs {
		ch <- p
	}
	close(ch)
}

func printProverbs(ch chan *proverb) {
	for p := range ch {
		fmt.Printf("%s\n", p.line)
		fmt.Printf("%s\n\n", formatMap(p.chars))
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
