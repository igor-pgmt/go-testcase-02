package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

// importFile импортирует указанный файл в граф и запускает расчёт карты путей
func (g *graph) importFile(fileName string) {
	g.fileName = fileName
	g.generateMap()
	g.findPaths()
}

// checkFile проверяет файл на корректность содержимого
func checkFile(fileName string) bool {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		matched, err := regexp.MatchString(`^[a-zA-Zа-яА-Я]* [a-zA-Zа-яА-Я]* \d*$`, scanner.Text())
		if err != nil {
			panic(err)
		}
		if !matched {
			log.Printf("Error in this line:\n%s\n", scanner.Text())
			return false
		}
	}
	return true
}

// getFileContent возвращает контент файла в виде массива строк
func (g *graph) getFileContent() (rows []string) {
	file, err := os.Open(g.fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}
	return
}
