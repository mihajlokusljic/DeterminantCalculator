package determinantCalc

import "sync"

type MinorCalcTask struct {
	Matrix        [][]float64
	ColumnIndexes []int
	ResultsHolder *float64
}

func (mct *MinorCalcTask) publishResult(result float64) {
	*mct.ResultsHolder = result
}

func calculateMinor(minorCalcTask MinorCalcTask) {
	result := minorCalcSerial(minorCalcTask.Matrix, 1, minorCalcTask.ColumnIndexes)
	minorCalcTask.publishResult(result)
}

func minorsCalculation(calculationTasks []MinorCalcTask, parentSemaphore *sync.WaitGroup) {
	for _, task := range calculationTasks {
		calculateMinor(task)
	}

	parentSemaphore.Done()
}

func DetParallel(matrix [][]float64, tasksNum int) float64 {

	n := len(matrix)
	columnIndexes := make([]int, n)
	for i := 0; i < n; i++ {
		columnIndexes[i] = i
	}

	var result float64 = 0.0
	minors := make([]float64, n)
	var sgn float64 = 1.0

	// in order to calculate the determinant the minors must be calculated first
	// each task calculates some of the minors
	tasks := make([][]MinorCalcTask, tasksNum)
	var wg sync.WaitGroup
	wg.Add(tasksNum)

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
