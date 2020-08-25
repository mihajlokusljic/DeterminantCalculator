from itertools import filterfalse
from multiprocessing import Process, Array

from determinanat_calc.util import measure_exec_time
from IO.matrix_reader import read_matrix
from determinanat_calc.serial_det_calc import minor_calc

class Minor_calc_task:
    """
    Holds information required to calculate a minor of a matrix, when the first
    row is removed. For calculation, the original matrix is accessed (only the rows
    and columns which are necessary to calculate the given minor). The result is written
    to a given sequence, at the given index.

    Attributes:
        matrix (list(list(float)): original matrix for which the minor should be calculated
        column_indexes (list(int)): indexes of columns in the original matrix containing the submatrix of given minor
        results_holder (multiprocessing.Array): sequence to write the minor in (used for final determinant calculation)
        result_idx (int): index of the given minor in the sequence that holds the result
    """

    def __init__(self, matrix, column_indexes, results_holder, result_idx):
        self.matrix = matrix
        self.column_indexes = column_indexes
        self.results_holder = results_holder
        self.result_idx = result_idx

    def publish_result(self, result):
        """
        Writes the given result to the sequence that holds the result (minors needed to calculate the determinant).

        Args:
            result (float): resulting minor value

        Return:
            None
        """
        self.results_holder[self.result_idx] = result


def serial_minor_calc(minor_calc_task):
    """
    Calculates a single minor of matrix using serial implementation. Since we are convolving over first row,
    all minors begin at second row (with index 1). The result is stored in a sequence containing the minors
    required to calculate the final determinant.

    Args:
        minor_calc_task (Minor_calc_task): contains information required to calculate the given minor

    Return:
        None
    """
    result, _ = minor_calc(minor_calc_task.matrix, 1, minor_calc_task.column_indexes)
    minor_calc_task.publish_result(result)

def minors_calculation(calculation_tasks):
    """
    Calculates all the minors assigned to a single process.

    Args:
        calculation_tasks (list(Minor_calc_task)): information required to calculate the minors and store the results

    Return:
        None
    """
    for task in calculation_tasks:
        serial_minor_calc(task)

@measure_exec_time
def det_parallel(matrix, tasks_num):
    """
    Calculates the determinant of given matrix using parallel implementation with given number of tasks.
    Task number parametrization is useful for scaling experiments. Let n be the order of the matrix
    and let p be the number of tasks. To calculate the determinant the values of n minors are required.
    The calculation of these n minors are assigned to p tasks using round Robin policy. After all tasks
    are finished calculating the required minors, the final determinant is calculated and returned.
    Since the methods is decorated with measure_exec_time, the execution time in milliseconds is also returned.

    Args:
        matrix (list(list(float)): matrix for which the determinant is calculated
        tasks_num (int): number of processes to be used in parallel calculation

    Return:
        value of the determinant and execution time in milliseconds ( (float, float) )
    """
    n = len(matrix)
    column_indexes = [i for i in range(n)]

    result = 0
    minors = Array('d', n)
    sgn = 1

    # in order to calculate the determinant the minors must be calculated first
    # each task calculates some of the minors
    child_processes = []
    tasks = [[] for i in range(tasks_num)]
    for idx, col in enumerate(column_indexes):
        # prepare the task holding information needed to calculate the given minor
        minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
        minor_calc_task = Minor_calc_task(matrix, minor_cols, minors, idx)
        # add this task to one of the processes, using round Robin policy
        p_idx = idx % tasks_num
        tasks[p_idx].append(minor_calc_task)

    # start tasks to calculate minors (first task set is executed by the current process)
    tasks_it = iter(tasks)
    first_task = next(tasks_it)

    for calc_tasks in tasks_it:
        p = Process(target=minors_calculation, args=(calc_tasks, ))
        child_processes.append(p)
        p.start()

    minors_calculation(first_task)

    # wait for other tasks to finish calculating assigned minors
    for p in child_processes:
        p.join()

    # calculate the determinant
    for j in range(n):
        result += sgn * matrix[0][column_indexes[j]] * minors[j]
        sgn *= -1

    return result


if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant, exec_time_ms = det_parallel(test_matrix, 5)
    print('det(matrix) =', determinant)
    print('Execution took: {} ms.'.format(exec_time_ms))