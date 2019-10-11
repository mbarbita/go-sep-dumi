package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// time conversions
// 15/03/2019 01:00:00;570389.25
const FormatTimeIn = "02/01/2006 15:04:05"

// Flags
type CmdArgs struct {
	dataFileIn string
}

var outFName string
var outFile *os.File
var outWriter *bufio.Writer

var cmdArgs = new(CmdArgs)

// Print line and character positions for easy slice ref
func printLineNr(st string) {
	fmt.Println("Line:", st)
	fmt.Println("Line Len:", len(st))
	for i := 0; i < len(st); i++ {
		fmt.Printf("%2v|", st[i:i+1])
	}
	fmt.Print("\n")
	for i := 0; i < len(st); i++ {
		fmt.Printf("%2v|", i)
	}
	fmt.Print("\n\n")
}

func convUTF16(s string) string {
	raw := []byte(s)
	// Read the file into a []byte:
	// raw, err := ioutil.ReadFile(filename)
	// if err != nil {
	// return nil, err
	// }

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		log.Println("utf16 conv error")
	}
	return string(decoded)
}

func init() {
	flag.StringVar(&cmdArgs.dataFileIn, "in", "src.csv", "filename")
	flag.Parse()

}

func main() {
	tstart := time.Now()
	fmt.Println("Started...")
	fmt.Println("input file:", cmdArgs.dataFileIn)
	fmt.Println()

	outFName = "export/" + cmdArgs.dataFileIn
	outFile, err := os.Create(outFName)
	if err != nil {
		log.Println("create file out error:")
		log.Fatal(err)
	}
	outWriter = bufio.NewWriter(outFile)
	defer outFile.Close()

	i := 1
	media := 0.
	count := 0
	sum := 0.

	// Open source csv file
	dataFileIn, err := os.Open(cmdArgs.dataFileIn)
	if err != nil {
		log.Println("open source file error:")
		log.Fatal(err)
	}
	defer dataFileIn.Close()

	// Read lines from file

	scanner := bufio.NewScanner(dataFileIn)

	//Loop over dataFileIn
	for scanner.Scan() {

		// Skip first line - table header
		if i == 1 {
			// fmt.Println("Skip line:", i, strings.TrimSpace(scanner.Text()))
			fmt.Println("skip line:", i)
			// fmt.Println(convUTF16(scanner.Text()))
			i++
			continue
		}

		if i == 2 {
			// fmt.Println("line:", i, convUTF16(scanner.Text()))
			// printLineNr(scanner.Text())
			// printLineNr(convUTF16(scanner.Text()))
			// t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
			// s, _, _ := transform.String(scanner.Text())
			// fmt.Println(s)
		}

		fields := strings.Split(scanner.Text(), ";")
		if len(fields) < 2 {
			break
		}
		// Process line into variables
		f1 := fields[0]
		// printLineNr(f1)
		f12 := ""
		for j := 1; j < len(f1); j += 2 {
			f12 += string(f1[j])
		}
		// fmt.Println("f12:", f12)
		f2 := fields[1]
		// printLineNr(f2)
		// f2 = strings.TrimSpace(f2)
		f22 := ""
		for j := 1; j < len(f2)-2; j += 2 {
			f22 += string(f2[j])
		}
		// fmt.Println("f22:", f22)
		// fmt.Println(f2)

		// Parse time
		timeIn, err := time.Parse(FormatTimeIn, f12)
		if err != nil {
			log.Println("time parse error:")
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Prepare/convert variables

		deb, err := strconv.ParseFloat(strings.TrimSpace(f22), 64)
		if err != nil {
			log.Println("deb parse error:")
			panic(err.Error())
		}
		sum += deb
		count++

		// Insert data into file
		if (timeIn.Minute() == 0) && (timeIn.Second() == 0) {
			media = sum / float64(count)
			// fmt.Println("i:", i, media, sum, count)
			fmt.Printf("i: %v time: %v media: %6.0f\n",
				i,
				timeIn.Format(FormatTimeIn), media)

			sum = 0
			count = 0
			fmt.Fprintln(outFile, fmt.Sprintf("%v;%6.0f",
				timeIn.Format(FormatTimeIn), media))
			// fmt.Println(i)
		}
		i++

		if i == 11 {
			break
		}

	} // Loop over dataFileIn end
	if err := scanner.Err(); err != nil {
		log.Println("input file scanner error:")
		log.Fatal(err)
	}

	outWriter.Flush()
	// Print summary

	fmt.Println()
	fmt.Println("START:", tstart.Format("15:04:05"))
	fmt.Println("END:  ", time.Now().Format("15:04:05"))
	fmt.Println("DURATION:", time.Since(tstart))
	fmt.Println("\a")

}
