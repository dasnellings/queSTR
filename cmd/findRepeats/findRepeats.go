package main

import (
	"flag"
	"fmt"
	"github.com/dasnellings/queSTR/seq"
	"log"
	"os"
	"strings"
)

type bed struct {
	Chr   string
	Start int
	End   int
	Name  string
}

func usage() {
	fmt.Print(
		"findRepeats - Find perfect repeats in a fasta file.\n" +
			"Usage:\n" +
			"findRepeats [options] -i input.fasta > repeats.bed\n\n")
	flag.PrintDefaults()
}

func main() {
	input := flag.String("i", "", "Input fasta file.")
	minUnitLen := flag.Int("minUnitSize", 2, "Minimum length of repeat unit to search for.")
	maxUnitLen := flag.Int("maxUnitSize", 6, "Maximum length of repeat unit to search for.")
	minUnitCount := flag.Int("minUnitCount", 5, "Minimum number of consecutive repeat units to be output.")
	flag.Parse()

	if *input == "" {
		usage()
		log.Fatal("ERROR: Must input a fasta file.")
	}

	findRepeats(*input, *minUnitLen, *maxUnitLen, *minUnitCount)
}

func findRepeats(input string, minUnitLen, maxUnitLen, minUnitCount int) {
	seqs := seq.Read(input)
	for name := range seqs {
		search(name, seqs[name], minUnitLen, maxUnitLen, minUnitCount)
	}
}

func search(chr, seq string, minUnitLen, maxUnitLen, minUnitCount int) {
	var idx, currUnitLen, consecutiveUnits int
	var motif string
	var b bed
	for idx < len(seq) { // walk across chromosome
		for currUnitLen = minUnitLen; currUnitLen <= maxUnitLen; currUnitLen++ { // check for repeats as we walk
			consecutiveUnits, motif = checkForRepeat(seq[idx:], currUnitLen)
			if consecutiveUnits > minUnitCount {
				break
			}
		}

		// if we found a repeat, write to output and move up, else increment idx
		if consecutiveUnits < minUnitCount {
			idx++
			continue
		}

		// if we got here, we found a repeat
		b.Chr = chr
		b.Start = idx
		b.End = idx + (consecutiveUnits * len(motif))
		b.Name = fmt.Sprintf("%dx%s", consecutiveUnits, strings.ToUpper(motif))

		output(b)

		idx += b.End - b.Start
	}
}

func checkForRepeat(seq string, unitLen int) (consecutiveUnits int, motif string) {
	if unitLen > len(seq) {
		return
	}
	motif = seq[:unitLen]
	if motif[0] == motif[1] {
		return
	}
	for (consecutiveUnits*unitLen)+unitLen < len(seq) {
		if seq[consecutiveUnits*unitLen:(consecutiveUnits*unitLen)+unitLen] == motif {
			consecutiveUnits++
		} else {
			return
		}
	}
	return
}

func output(b bed) {
	_, err := fmt.Fprintf(os.Stdout, "%s\t%d\t%d\t%s\n", b.Chr, b.Start, b.End, b.Name)
	if err != nil {
		log.Panic(err)
	}
}
