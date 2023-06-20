package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		file *os.File
		f    *os.File
		err  error
	)
	filePath := "../traces/Financial1.spc"
	if file, err = os.Open(filePath); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err = os.OpenFile("../traces_out/fin1_not_convert", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")
		//if row[3] == "Write" {
		//	row[3] = "W"
		//} else {
		//	row[3] = "R"
		//}
		if _, err := f.WriteString(row[1] + "," + row[3] + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("Done")
}
