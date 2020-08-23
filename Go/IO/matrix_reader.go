package IO

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// ReadMatrix loads and returns a regular matrix from given text file. The file must follow the convention:
// first row contains a non-negative integer number (n) which specifies the order of the matrix.
// The next n rows contain exactly n real numbers, separated by one or more space characters.
//
// Args:
// 	matrixFilePath: path of the file containing the matrix
//
// Return:
// 	matrix loaded into memory and ready for use
func ReadMatrix(matrixFilePath string) [][]float64 {

	f, err := os.Open(matrixFilePath)

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
