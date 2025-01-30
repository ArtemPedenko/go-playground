package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func recDir(path string, prefixes []bool, printFiles bool) (string, error) {
	var returnString string

	dir, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	for fileIndex, file := range dir {
		fullPath := path + "/" + file.Name()

		// Формируем отступы для текущего уровня
		var prefixStr strings.Builder
		for _, p := range prefixes {
			if p {
				prefixStr.WriteString("│ ")
			} else {
				prefixStr.WriteString("  ")
			}
		}

		// Определяем символ ветвления
		var branching string

		if fileIndex == len(dir)-1 {
			branching = "└───"
		} else {
			branching = "├───"
		}

		if file.IsDir() {
			returnString += prefixStr.String() + branching + file.Name() + "\n"

			subDirString, err := recDir(fullPath, append(prefixes, fileIndex != len(dir)-1), printFiles)
			if err != nil {
				return "", err
			}

			returnString += subDirString
		} else if printFiles {
			stat, err := os.Stat(fullPath)
			if err != nil {
				return "", err
			}

			// Размер файла
			size := strconv.FormatInt(stat.Size(), 10)
			sizeString := "(empty)"
			if stat.Size() > 0 {
				sizeString = "(" + size + "b)"
			}

			returnString += prefixStr.String() + branching + file.Name() + " " + sizeString + "\n"
		}
	}

	return returnString, nil
}

// вывод файлов и папок в директории
// go run main go <путь>			(папки)
// go run main.go <путь> -f   (с файлами)
// go run main.go . -f
func main() {
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	res, err := recDir(path, []bool{}, printFiles)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(res)
}
