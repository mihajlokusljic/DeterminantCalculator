from itertools import filterfalse
from IO.matrix_reader import read_matrix

from determinanat_calc.util import measure_exec_time

def minor_calc(matrix, begin_row_index, column_indexes):
    """
    Calculates an arbitrary minor of given matrix. The original matrix is accessed
    for calculation (only the necessary rows and columns).

    Args:
        matrix (list(list(float))): matrix containing the submatrix of the minor
        begin_row_index (int): index of the first row of the submatrix for the minor, in the original matrix
        column_indexes (list(int)): indexes of columns of the submatrix for the minor, in the original matrix

    Return:
        value of the given minor (float)
    """

    # order of submatrix
    n = len(column_indexes)

    # if order is one, return the only element of the submatrix as the minor
    if n == 1:
        return matrix[begin_row_index][column_indexes[0]]

    # initialize the values to calculate the determinant
    result = 0
    minors = [0 for i in range(n)]
    sgn = 1

    # convolve over first row
    for idx, col in enumerate(column_indexes):
        minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
        minors[idx] = minor_calc(matrix, begin_row_index + 1, minor_cols)  # calculate required minors

    for j in range(n):
        result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
        sgn *= -1

    return result

@measure_exec_time
def det_serial(matrix):
    """
    Calculates and returns the determinant of given matrix using serial implementation.
    Since it is decorated with measure_exec_time, the execution time in milliseonds is also returned.

    Args:
        matrix (list(list(float)): matrix for which the determinant is calculated.

    Return:
        value of the determinant and execution time in milliseconds ( (float, float) )
    """
    n = len(matrix)
    cols = list([i for i in range(n)])
    return minor_calc(matrix, 0, cols)

if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant, exec_time_ms = det_serial(test_matrix)
    print('det(matrix) =', determinant)
    print('Execution time was: {} ms.'.format(exec_time_ms))





