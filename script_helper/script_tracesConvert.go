package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	searchDir := "../traces/MSR-Cambridge/"
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".csv") {
			writeToFile(path)
		}
		return nil
	})
	if e != nil {
		return //fmt.Errorf("filepath.Walk() returned %v\n", e)
	}
}

func writeToFile(path string) {
	var (
		file *os.File
		err  error
		f    *os.File
	)

	file, err = os.OpenFile("../traces_out/msr_src1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if f, err = os.Open(path); err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")
		if row[3] == "Write" {
			row[3] = "W"
		} else {
			row[3] = "R"
		}
		if _, err := file.WriteString(row[4] + "," + row[3] + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("Done")
}
