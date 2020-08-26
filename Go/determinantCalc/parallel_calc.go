package determinantCalc

import "sync"

// MinorCalcTask holds information required to calculate a minor of a matrix, when the first
// row is removed. For calculation, the original matrix is accessed (only the rows
// and columns which are necessary to calculate the given minor).
//
// Attributes:
// 	Matrix: original matrix for which the minor should be calculated
// 	ColumnIndexes: indexes of columns in the original matrix containing the submatrix of given minor
// 	ResultsHolder: points to where the resulting value shoud be written/stored
type MinorCalcTask struct {
	Matrix        [][]float64
	ColumnIndexes []int
	ResultsHolder *float64
}

// publishResult method writes the given result to the sequence that holds the result (minors needed to calculate the determinant).
//
// Args:
// 	result: resulting minor value
//
// Return:
// 	nil
func (mct *MinorCalcTask) publishResult(result float64) {
	*mct.ResultsHolder = result
}

// Calculates a single minor of matrix using serial implementation. Since we are using Laplace expansion over first row,
// all minors begin at second row (with index 1). The result is stored in a sequence containing the minors
// required to calculate the final determinant.
//
// Args:
//     minorCalcTask: contains information required to calculate the given minor
//
// Return:
//     nil
func calculateMinor(minorCalcTask MinorCalcTask) {
	result, _ := minorCalcSerial(minorCalcTask.Matrix, 1, minorCalcTask.ColumnIndexes, false)
	minorCalcTask.publishResult(result)
}

func minorsCalculation(calculationTasks []MinorCalcTask, parentSemaphore *sync.WaitGroup) {
	for _, task := range calculationTasks {
		calculateMinor(task)
	}

	parentSemaphore.Done()
}

// DetParallel calculates the determinant of given matrix using parallel implementation with given number of tasks.
// Task number parametrization is useful for scaling experiments. Let n be the order of the matrix
// and let p be the number of tasks. To calculate the determinant the values of n minors are required.
// The calculation of these n minors are assigned to p tasks using round Robin policy. After all tasks
// are finished calculating the required minors, the final determinant is calculated and returned.
//
// Args:
// 	matrix: matrix for which the determinant is calculated
// 	tasksNum: number of processes to be used in parallel calculation
//
// Return:
// 	value of the determinant
func DetParallel(matrix [][]float64, tasksNum int) float64 {

	// the starting minor is the entire matrix, with all of its columns
	n := len(matrix)
	columnIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		columnIndexes[i] = i
	}

	// initialize the result and allocate memory to hold values of minors
	var result float64 = 0.0
	minors := make([]float64, n)
	var sgn float64 = 1.0

	// in order to calculate the determinant the minors must be calculated first,
	// each task calculates some of the minors

	// for each process holds assigned minor calculation tasks
	tasks := make([][]MinorCalcTask, tasksNum)
	// semaphore used to sync tasks
	var wg sync.WaitGroup
	wg.Add(tasksNum)

	// for each requored minor form a task and assign it to one of the processes
	for idx, col := range columnIndexes {
		// prepare the task holding information needed to calculate the given minor
		minorCols := removeColumn(columnIndexes, col)
		minorCalcTask := MinorCalcTask{
			Matrix:        matrix,
			ColumnIndexes: minorCols,
			ResultsHolder: &minors[idx],
		}

		// add this task to one of the processes, using round Robin policy
		pIdx := idx % tasksNum
		tasks[pIdx] = append(tasks[pIdx], minorCalcTask)
	}

	// start tasks to calculate minors (first task set is executed by the current process)
	firstTask := tasks[0]

	for _, calcTasks := range tasks[1:] {
		go minorsCalculation(calcTasks, &wg)
	}

	minorsCalculation(firstTask, &wg)

	// wait for other tasks to finish calculating assigned minors
	wg.Wait()

	// calculate the determinant
	for j := 0; j < n; j++ {
		result += sgn * matrix[0][columnIndexes[j]] * minors[j]
		sgn *= -1
	}

	return result
}
