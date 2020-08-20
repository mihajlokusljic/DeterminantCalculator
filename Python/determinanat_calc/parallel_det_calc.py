from itertools import filterfalse
from multiprocessing import Process, Array

from determinanat_calc.util import measure_exec_time
from IO.matrix_reader import read_matrix
from determinanat_calc.serial_det_calc import minor_calc

class Minor_calc_task:

    def __init__(self, matrix, column_indexes, results_holder, result_idx):
        self.matrix = matrix
        self.column_indexes = column_indexes
        self.results_holder = results_holder
        self.result_idx = result_idx

    def publish_result(self, result):
        self.results_holder[self.result_idx] = result


def serial_minor_calc(minor_calc_task):
    # since we are convolving over first row, all minors begin at second row (with index 1)
    result = minor_calc(minor_calc_task.matrix, 1, minor_calc_task.column_indexes)
    minor_calc_task.publish_result(result)

def minors_calculation(calculation_tasks):
    for task in calculation_tasks:
        serial_minor_calc(task)

@measure_exec_time
def det_parallel(matrix, tasks_num):
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