package IO

import (
	"fmt"
	"log"
	"os"
)

type ExecutionResult struct {
	MatrixOrder int
	ExecTimeMs  int64
	Serial      bool
}

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
