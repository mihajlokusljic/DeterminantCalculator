from itertools import filterfalse
from multiprocessing import Process, Array

from IO.matrix_reader import read_matrix
from determinanat_calc.serial_det_calc import det_serial

CUTOFF = 8


def minor_calc(matrix, begin_row_index, column_indexes, results_holder, result_idx):
    n = len(column_indexes) # order of matrix

    if n <= CUTOFF:
        results_holder[result_idx] = det_serial(matrix, begin_row_index, column_indexes)
    else:
        result = 0
        # minors = [0 for i in range(n)]
        minors = Array('d', n)
        sgn = 1
        child_processes = []

        # razvoj po prvoj vrsti
        for idx, col in enumerate(column_indexes):
            minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
            minor_calc_task = Process(target=minor_calc, args=(matrix, begin_row_index + 1, minor_cols, minors, idx))
            child_processes.append(minor_calc_task)
            minor_calc_task.start()
            # minor_calc(matrix, begin_row_index + 1, minor_cols, minors, idx) # calculate required minors

        # wait for child tasks to finish calculating required minors
        for p in child_processes:
            p.join()

        # calculate the given minor
        for j in range(n):
            result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
            sgn *= -1

        results_holder[result_idx] = result

def det_parallel(matrix):
    result_holder = [0]
    n = len(matrix)
    column_indexes = [i for i in range(n)]
    minor_calc(matrix, 0, column_indexes, result_holder, 0)
    return result_holder[0]

if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant = det_parallel(test_matrix)
    print('det(matrix) =', determinant)