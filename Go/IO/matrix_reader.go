package IO

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ReadMatrix(matrix_file_path string) [][]float64 {

	f, err := os.Open(matrix_file_path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	order := 0
	scanner.Scan()
	order, err = strconv.Atoi(scanner.Text())

	matrix := make([][]float64, order)
	for i := 0; i < order; i++ {
		matrix[i] = make([]float64, order)
	}

	var currentVal float64

	for i := 0; i < order; i++ {
		for j := 0; j < order; j++ {
			scanner.Scan()
			currentVal, err = strconv.ParseFloat(scanner.Text(), 64)
			matrix[i][j] = currentVal
		}
	}

	return matrix

}
