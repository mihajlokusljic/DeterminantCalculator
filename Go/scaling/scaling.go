package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/mihajlokusljic/DeterminantCalculator/IO"
	"gitlab.com/mihajlokusljic/DeterminantCalculator/determinantCalc"
)

const matrixPathTemplate = "../../test_data/matrica%dx%d.txt"
const resultsBasePath = "../../results"
const strongScalingMatrixOrder = 10

var availableMatrixOrders = [...]int{3, 5, 8, 9, 10, 11}

type ScalingResult struct {
	MatrixOrder int
	TasksNum    int
	ExecTimeMs  int64
}

func writeScalingResults(resultsFilePath string, scalingResults []ScalingResult) {
	resultsFile, err := os.OpenFile(resultsFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to access results file.\n", err)
	}
	defer resultsFile.Close()

	resultsFile.WriteString("matrix_order,number_of_tasks,exec_time_ms\n")
	for _, res := range scalingResults {
		resultsFile.WriteString(fmt.Sprintf("%d,%d,%d\n",
			res.MatrixOrder, res.TasksNum, res.ExecTimeMs))
	}
}

func loadTestMatrix(matrixOrder int) [][]float64 {
	var matrixFilePath string = fmt.Sprintf(matrixPathTemplate, matrixOrder, matrixOrder)
	var matrix [][]float64 = IO.ReadMatrix(matrixFilePath)
	fmt.Printf("Successfully loaded matrix of order %d. Calculating determinant(s)...\n", len(matrix))
	return matrix
}

func strongScalingExperiment() {
	fmt.Println("Strong scaling experiment:\n")
	var matrix [][]float64 = loadTestMatrix(strongScalingMatrixOrder)
	var results []ScalingResult

	for tasksNum := 1; tasksNum <= strongScalingMatrixOrder; tasksNum++ {
		start := time.Now()
		determinantCalc.DetParallel(matrix, tasksNum)
		elapsed := time.Since(start)
		duration := elapsed.Milliseconds()
		scalingResult := ScalingResult{
			ExecTimeMs:  duration,
			MatrixOrder: strongScalingMatrixOrder,
			TasksNum:    tasksNum,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n",
			strongScalingMatrixOrder, tasksNum, duration)
	}

	writeScalingResults(fmt.Sprintf("%s/strong_scaling_results_go.csv", resultsBasePath), results)
}

func weakScalingExperiment() {
	fmt.Println("\nWeak scaling experiment:\n")
	var results []ScalingResult

	for _, n := range availableMatrixOrders {
		matrix := loadTestMatrix(n)

		start := time.Now()
		determinantCalc.DetParallel(matrix, n)
		elapsed := time.Since(start)
		duration := elapsed.Milliseconds()
		scalingResult := ScalingResult{
			ExecTimeMs:  duration,
			MatrixOrder: n,
			TasksNum:    n,
		}
		results = append(results, scalingResult)

		fmt.Printf("Parallel determinant calculation of matrix of order %d with %d tasks took %d ms.\n", n, n, duration)
	}

	writeScalingResults(fmt.Sprintf("%s/weak_scaling_results_go.csv", resultsBasePath), results)
}

func main() {
	strongScalingExperiment()
	weakScalingExperiment()
}
