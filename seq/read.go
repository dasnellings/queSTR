package seq

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

// Read reads a fasta file to a map of sequence strings keyed by the record name.
func Read(filename string) map[string]string {
	ans := make(map[string]string)
	buf := new(bytes.Buffer)
	file, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)

	var line, currName string
	for sc.Scan() {
		line = sc.Text()
		if strings.HasPrefix(line, ">") {
			if buf.Len() > 0 {
				ans[currName] = buf.String()
				buf.Reset()
			}
			currName = line[1:]
			continue
		}

		buf.WriteString(line)
	}

	if err = sc.Err(); err != nil {
		log.Panic(err)
	}

	if buf.Len() > 0 {
		ans[currName] = buf.String()
	}

	return ans
}
