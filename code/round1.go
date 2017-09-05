package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

/*
	round1.go
		by Gem

	Initial round1 will pull the dataset from: http://icon.shef.ac.uk/Moby/mpos.html
	and normalized it into a simple CSV file

	CSV file format:
		Column1: Word
		Column2: Definition (if doesn't exit - Nil)
		Column3: Part of Speech (if doesn't exist - Nil) (Separated by -)
		Column4: Part of Speech Symbol (if doesn't exist - Nil) (Separated by -)
*/

const (
	directory = `C:\Users\mshannon\Desktop\GoSpace\src\github.com\gemmyboy\enlang-dataset\raw\mpos\mobyposi.txt`
	csvfile   = `C:\Users\mshannon\Desktop\GoSpace\src\github.com\gemmyboy\enlang-dataset\enlang-dataset.csv`
)

func main() {

	moby, err := os.Open(directory)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting 1st Pass", moby.Name())

	wordMap := make(map[string][]string) //map[word]pos

	//First pass is to strip data from file
	reader := bufio.NewReader(moby)
	for {
		line, errr := reader.ReadString('\n')
		if errr != nil {
			fmt.Println(errr)
			break
		}
		fmt.Println(line)

		columns := strings.Split(line, string(215)) //POS is split by ASCII 215
		wordMap[columns[0]] = []string{columns[1], ""}
		fmt.Println("hey")
	}

	moby.Close()

	fmt.Println("Starting 2nd Pass")

	//Second pass is to handle POS and create the 3rd & 4th column
	for key, value := range wordMap {
		col3 := bytes.NewBufferString("")
		col4 := bytes.NewBufferString("")
		for _, c := range value[0][:] {
			switch string(c) {
			case "N":
				col3.WriteString("Noun-")
				col4.WriteString("N-")
			case "p":
				col3.WriteString("Plural-")
				col4.WriteString("p-")
			case "h":
				col3.WriteString("Noun Phrase-")
				col4.WriteString("h-")
			case "V":
				col3.WriteString("Verb Participle-")
				col4.WriteString("V-")
			case "t":
				col3.WriteString("Verb Transitive-")
				col4.WriteString("t-")
			case "i":
				col3.WriteString("Verb Intransitive-")
				col4.WriteString("i-")
			case "A":
				col3.WriteString("Adjective-")
				col4.WriteString("A-")
			case "v":
				col3.WriteString("Adverb-")
				col4.WriteString("v-")
			case "C":
				col3.WriteString("Conjunction-")
				col4.WriteString("C-")
			case "P":
				col3.WriteString("Preposition-")
				col4.WriteString("P-")
			case "!":
				col3.WriteString("Interjection-")
				col4.WriteString("!-")
			case "r":
				col3.WriteString("Pronoun-")
				col4.WriteString("r-")
			case "D":
				col3.WriteString("Definite Article-")
				col4.WriteString("D-")
			case "I":
				col3.WriteString("Indefinite Article-")
				col4.WriteString("I-")
			case "o":
				col3.WriteString("Nominative-")
				col4.WriteString("o-")
			}
		}

		tCol3 := col3.String()
		tCol4 := col4.String()
		wordMap[key] = []string{tCol3[:(len(tCol3) - 1)], tCol4[:(len(tCol4) - 1)]}
	}

	fmt.Println("Starting 3rd Pass")

	//Third pass is to generate the file
	csvFile, er := os.Create(csvfile)
	if er != nil {
		panic(er)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)

	for word, data := range wordMap {
		record := []string{word, "nil", data[0], data[1]}
		w.Write(record)
	}

	w.Flush()
} //End main()
