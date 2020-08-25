from itertools import filterfalse
from IO.matrix_reader import read_matrix
import time

from determinanat_calc.util import measure_exec_time

def minor_calc(matrix, begin_row_index, column_indexes, measure_parallel_code = False):
    """
    Calculates an arbitrary minor of given matrix. The original matrix is accessed
    for calculation (only the necessary rows and columns). Returns two values:
        1. the value of the minor
        2. time in milliseconds spent executing code that can be parallelized if argument
        measure_parallel_code is True, and zero otherwise (it will not be measured)

    Args:
        matrix (list(list(float))): matrix containing the submatrix of the minor
        begin_row_index (int): index of the first row of the submatrix for the minor, in the original matrix
        column_indexes (list(int)): indexes of columns of the submatrix for the minor, in the original matrix
        measure_parallel_code (bool): indicates whether the execution time for code that can be parallelized
        should be measured, it is False by default for better performance

    Return:
        value of the given minor, time spent executing code that can be parallelized ( (float, float) )
    """

    # order of submatrix
    n = len(column_indexes)

    # if order is one, return the only element of the submatrix as the minor
    if n == 1:
        return matrix[begin_row_index][column_indexes[0]], 0.0

    # initialize the values to calculate the determinant
    result = 0
    minors = [0 for i in range(n)]
    sgn = 1

    # expansion over first row
    parallel_code_exec_time = 0
    for idx, col in enumerate(column_indexes):
        minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
        # calculate required minors (this can be done in parallel)
        if measure_parallel_code:
            start_time = time.time()
            minor, _ = minor_calc(matrix, begin_row_index + 1, minor_cols)
            end_time = time.time()
            minor_calc_time = (end_time - start_time) * 1000
            minors[idx] = minor
            parallel_code_exec_time += minor_calc_time
        else:
            minor, _ = minor_calc(matrix, begin_row_index + 1, minor_cols)
            minors[idx] = minor


    for j in range(n):
        result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
        sgn *= -1

    return result, parallel_code_exec_time

@measure_exec_time
def det_serial(matrix):
    """
    Calculates and returns the determinant of given matrix using serial implementation, as well
    as time in milliseconds, spent executing code that can be parallelized.
    Since it is decorated with measure_exec_time, the execution time in milliseonds is also returned.

    Args:
        matrix (list(list(float)): matrix for which the determinant is calculated.

    Return:
        value of the determinant, time spent executing code that can be parallelized
        and total execution time in milliseconds ( (float, float, float) )
    """
    n = len(matrix)
    cols = list([i for i in range(n)])
    determinant, parallel_code_exec_time = minor_calc(matrix, 0, cols, measure_parallel_code=True)
    return determinant, parallel_code_exec_time

if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant, parallel_code_exec_time, exec_time_ms = det_serial(test_matrix)
    print('det(matrix) =', determinant)
    print('Execution time was: {} ms.'.format(exec_time_ms))





