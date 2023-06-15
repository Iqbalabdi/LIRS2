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

	searchDir := "./traces"
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".blkparse") {
			writeToFile(path)
		}
		return nil
	})
	if e != nil {
		panic(e)
	}
}

func writeToFile(path string) {
	var (
		file *os.File
		err  error
		f    *os.File
	)

	file, err = os.OpenFile("./traces_out/fiu.web", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		row := strings.Split(scanner.Text(), " ")
		if _, err := file.WriteString(row[3] + "," + row[5] + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("Done")

}
