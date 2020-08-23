package IO

import (
	"fmt"
	"log"
	"os"
)

// ExecutionResult contains information about execution times for determinant calculations, which
// are stored in the results file.
//
// Attributes:
// 	MatrixOrder: order of the matrix for which the determinant was calculated
// 	ExecTimeMs: duration of the calculation process in milliseconds
// 	Serial (bool): if true the implementation is serial, otherwise implementation is parallel
type ExecutionResult struct {
	MatrixOrder int
	ExecTimeMs  int64
	Serial      bool
}

// WriteExecutionResults writes information about execution times of calculating determinants of
// different matrices, using serial and parallel implementation, to a specified CSV file.
// The CSV file header contains: n (matrix order), exec_time_ms (execution time), serial (determines
// whether the implementation is serial or parallel), implementation (determines the programming language
// of used implementation). If the file does not exist it will be created and initialized,
// otherwise the new results will be added to the end of existing file.
//
// Args:
// 	results: information about determinant calculation processes for different matrices and implementations
// 	resultFilePath: path of the CSV file to hold the information
//
// Return:
// 	nil
func WriteExecutionResults(results []ExecutionResult, resultFilePath string) {
	resultsFile, err := os.OpenFile(resultFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to access results file.\n", err)
	}
	defer resultsFile.Close()

	// if file is empty add the header
	if info, _ := resultsFile.Stat(); info.Size() == 0 {
		resultsFile.WriteString("n,exec_time_ms,serial,implementation\n")
	}

	// write execution results
	for _, res := range results {
		resultsFile.WriteString(fmt.Sprintf("%d,%d,%t,go\n",
			res.MatrixOrder, res.ExecTimeMs, res.Serial))
	}
}
