package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Upscale struct {
	amd    string
	nvidia string
	intel  string
	misc   string
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

	// re, err = regexp.Compile(`\|-\|\[\[([A-Za-z\-,\.&0-9™:'\(\)\|\! 	]+)\]\](<ref>.+<\/ref>)\|(.+)\|(.+)\|(.+)\n?`)
	re, err = regexp.Compile(`\|-\|\[\[([A-Za-z\-,\.&0-9™:–'\(\)\|\! 	]+)\]\](<ref>.+<\/ref>)\|(.+)\|(.+)\|(.+)\n?`)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(inputFile)

	const absent = "❌"

	for scanner.Scan() {
		line += 1

		if !re.Match([]byte(scanner.Text())) {
			fmt.Printf("Skipping... `%s`\n\n", scanner.Text())
			continue
		}

		res := re.FindStringSubmatch(scanner.Text())
		var amd, nvidia, intel, misc string
		technology := res[3]

		if strings.Contains(technology, "FSR 1") {
			amd += "FSR 1 / "
		}
		if strings.Contains(technology, "FSR 2.1") {
			amd += "FSR 2.1 / "
		} else if strings.Contains(technology, "FSR 2") {
			amd += "FSR 2.0 / "
		}

		if len(amd) == 0 {
			amd = absent
		}

		amd = strings.TrimSuffix(amd, " / ")

		if strings.Contains(technology, "NIS") {
			nvidia += "NIS / "
		}

		if strings.Contains(technology, "DLSS 1") {
			nvidia += "DLSS 1 / "
		} else if strings.Contains(technology, "DLSS 2") {
			nvidia += "DLSS 2 / "
		}

		if len(nvidia) == 0 {
			nvidia = absent
		}

		nvidia = strings.TrimSuffix(nvidia, " / ")

		if strings.Contains(technology, "XeSS") {
			intel += "XeSS"
		} else {
			intel += absent
		}

		if strings.Contains(technology, "TSR") {
			misc += "TSR / "
		}
		if strings.Contains(technology, "TAAU") {
			misc += "TAAU / "
		}

		if len(misc) == 0 {
			misc = absent
		}

		misc = strings.TrimSuffix(misc, " / ")

		var release = strings.Replace(fmt.Sprintf(res[4], res[5]), "EXTRA string=", "", 1)
		release = strings.Replace(release, "%!", ` `, 1)
		release = strings.Replace(release, `colspan="2"`, ``, 1)
		release = strings.Replace(release, `}}`, ``, 1)

		games = append(games, Game{
			name: res[1],
			upscale: Upscale{
				amd:    amd,
				nvidia: nvidia,
				intel:  intel,
				misc:   misc,
			},
			release: release,
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

		outputFile.WriteString(fmt.Sprint("{{User:Mine18/Templates/High Fidelity Upscaling/row|",
			"[[", game.name, "]]", "|",
			game.release, "|",
			game.upscale.amd, "|",
			game.upscale.nvidia, "|",
			game.upscale.intel, "|",
			game.upscale.misc, "|",
			game.notes, "}}\n"))
	}
}
