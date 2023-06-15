package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func main() {
	var (
		file *os.File
		f    *os.File
		err  error
	)
	charset := "RW"
	filePath := "./traces_out/scan.txt"
	if file, err = os.Open(filePath); err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if f, err = os.Create("./traces_out/New_scan.txt"); err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if _, err := f.WriteString(scanner.Text() + "," + string(charset[rand.Intn(len(charset))]) + "\n"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("Done")
}
