from itertools import filterfalse
from IO.matrix_reader import read_matrix

from determinanat_calc.util import measure_exec_time

def minor_calc(matrix, begin_row_index, column_indexes):
    n = len(column_indexes) # order of submatrix

    if n == 1:
        return matrix[begin_row_index][column_indexes[0]]

    result = 0
    minors = [0 for i in range(n)]
    sgn = 1

    # razvoj po prvoj vrsti
    for idx, col in enumerate(column_indexes):
        minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
        minors[idx] = minor_calc(matrix, begin_row_index + 1, minor_cols)  # calculate required minors

    for j in range(n):
        result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
        sgn *= -1

    return result

@measure_exec_time
def det_serial(matrix):
    n = len(matrix)
    cols = list([i for i in range(n)])
    return minor_calc(matrix, 0, cols)

if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant, exec_time_ms = det_serial(test_matrix)
    print('det(matrix) =', determinant)
    print('Execution time was: {} ms.'.format(exec_time_ms))





