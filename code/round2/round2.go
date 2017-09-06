package main

import "os"
import "io/ioutil"
import "fmt"
import "encoding/csv"
import "strings"

/*
	round2.go
		by Gem

	Found a reasonably suitable data file containing definitions for 176023 words.
	https://sourceforge.net/p/mysqlenglishdictionary/code/ci/master/tarball
	The MySql file wouldn't import in SSMS so I'm just going to strip them from the file manually.

	This assumes that there exists a file generated by round1.

	CSV file generated in round2 matches definitions from the MySql file to words in round 1.
*/

const (
	round1CSVPath = `C:\Users\mshannon\Desktop\GoSpace\src\github.com\gemmyboy\enlang-dataset\enlang-dataset-1.csv`
	round2CSVPath = `C:\Users\mshannon\Desktop\GoSpace\src\github.com\gemmyboy\enlang-dataset\enlang-dataset-2.csv`
	rawDictPath   = `C:\Users\mshannon\Desktop\GoSpace\src\github.com\gemmyboy\enlang-dataset\raw\mysqldict\dictionaryStudyTool.sql`
)

func main() {

	//First Parse the mysql file. It's been slightly cleaned for easier parsing.
	rawDict, er := os.Open(rawDictPath)
	if er != nil {
		panic(er)
	}

	data, err := ioutil.ReadAll(rawDict)
	if err != nil {
		panic(err)
	}
	rawDict.Close()

	fmt.Println("Starting 1st Pass: Stripping data from raw dictionary")

	//First Pass is to strip data into a 2D array
	dictData := make(map[string]string, 177000)
	current := []string{"", "", ""}
	iter := 0

	for i := 0; i < len(data); {
		if i+3 > len(data) {
			break
		}
		if string(data[i])+string(data[1]) == "('" {
			current = []string{"", "", ""}
			iter = 0
			i += 2
		} else if string(data[i])+string(data[i+1])+string(data[i+2]) == "','" {
			iter++
			i += 3
		} else if string(data[i])+string(data[i+1])+string(data[i+2]) == "')," {
			dictData[strings.ToLower(current[0])] = current[2]
			i += 3
			continue
		} else if string(data[i])+string(data[i+1]) == "')" {
			dictData[strings.ToLower(current[0])] = current[2]
			i += 2
			continue
		}

		if iter > 2 {
			fmt.Println(iter, i, current)
		}

		current[iter] = current[iter] + string(data[i])
		i++
	} //End for first

	fmt.Println("Starting 2nd Pass: Loading Round1 CSV")

	//Second pass is to load up the CSV file from Round1
	csv1, errr := os.Open(round1CSVPath)
	if errr != nil {
		panic(errr)
	}

	readCSV1 := csv.NewReader(csv1)
	csv1Data, errrr := readCSV1.ReadAll()
	if errrr != nil {
		panic(errrr)
	}
	csv1.Close()

	fmt.Println("Starting 3rd Pass: Matching data ")

	//Third pass is to match the data against each other
	for i, entry := range csv1Data {
		if def, ok := dictData[entry[0]]; ok {
			csv1Data[i][1] = def
		}
	}

	fmt.Println("Starting Final Pass: CSV Generation")

	//Final pass is to generated 2nd round CSV file
	csv2File, errrrr := os.Create(round2CSVPath)
	if errrrr != nil {
		panic(errrrr)
	}
	defer csv2File.Close()

	w := csv.NewWriter(csv2File)
	w.WriteAll(csv1Data)
	w.Flush()

	fmt.Println("Done")
} //End main()
