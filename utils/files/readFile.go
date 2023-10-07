package files

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func ReadTXT(path string) (string, error) {
	var res string
	op, err := os.Open(path)
	if err != nil {
		log.Println("Error while reading the txt file ")
		return "", err
	}
	defer op.Close()

	scanner := bufio.NewScanner(op)

	for scanner.Scan() {
		line := scanner.Text()
		res += line + "\n"
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
	return res, nil
}
