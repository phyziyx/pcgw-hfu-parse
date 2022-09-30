package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Upscale struct {
	fsr1  string
	fsr2  string
	dlss1 string
	dlss2 string
	xess  string
}

type Game struct {
	name    string
	upscale Upscale
	release string
	notes   string
}

func main() {
	inputFile, err := os.OpenFile("parse.txt", os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	var line = 0
	var re *regexp.Regexp
	var games = make([]Game, 0)
	var dlss1, dlss2, fsr1, fsr2, xess = "❌", "❌", "❌", "❌", "❌"

	re, err = regexp.Compile(`\|-\|\[\[([A-Za-z\-,\.&0-9™:'\(\)\|\! 	]+)\]\](<ref>.+<\/ref>)\|(.+)\|(.+)\|(.+)\n?`)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		line += 1

		if !re.Match([]byte(scanner.Text())) {
			fmt.Printf("Skipping... `%s`\n\n", scanner.Text())
			continue
		}

		res := re.FindStringSubmatch(scanner.Text())

		if strings.Contains(res[3], "FSR 1") {
			fsr1 = "✔"
		}

		if strings.Contains(res[3], "FSR 2") {
			fsr2 = "✔"
		}

		if strings.Contains(res[3], "DLSS 1") {
			dlss1 = "✔"
		}

		if strings.Contains(res[3], "DLSS 2") {
			dlss2 = "✔"
		}

		if strings.Contains(res[3], "xess") {
			xess = "✔"
		}

		games = append(games, Game{
			name: res[1],
			upscale: Upscale{
				fsr1:  fsr1,
				fsr2:  fsr2,
				dlss1: dlss1,
				dlss2: dlss2,
				xess:  xess,
			},
			release: fmt.Sprintf(res[4], res[5]),
			notes:   res[2],
		})
	}

	fmt.Println("Total lines:", line)
	fmt.Println("Outputting lines:", len(games))

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	for i := 0; i < len(games); i++ {
		var game = games[i]

		outputFile.WriteString(fmt.Sprint("|", game.name, "|", game.release, "|", game.upscale.fsr1, "/", game.upscale.fsr2, "|",
			game.upscale.dlss1, "/", game.upscale.dlss2, "|", game.upscale.xess, "|", game.notes, "|\n"))
	}
}
