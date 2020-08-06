from itertools import filterfalse
from IO.matrix_reader import read_matrix

def minor_calc(matrix, begin_row_index, column_indexes, results_holder, result_idx):
    n = len(column_indexes) # order of matrix

    if n == 1:
        results_holder[result_idx] = matrix[begin_row_index][column_indexes[0]]
    else:
        result = 0
        minors = [0 for i in range(n)]
        sgn = 1

        # razvoj po prvoj vrsti
        for idx, col in enumerate(column_indexes):
            minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
            minor_calc(matrix, begin_row_index + 1, minor_cols, minors, idx) # calculate required minors

        for j in range(n):
            result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
            sgn *= -1

        results_holder[result_idx] = result

def det_serial(matrix, begin_row_index = 0, column_indexes = []):
    n = len(column_indexes) # order of submatrix
    if begin_row_index == 0:
        n = len(matrix)
        column_indexes = list(range(n))

    if n == 1:
        return matrix[begin_row_index][column_indexes[0]]

    result = 0
    minors = [0 for i in range(n)]
    sgn = 1

    # razvoj po prvoj vrsti
    for idx, col in enumerate(column_indexes):
        minor_cols = list(filterfalse(lambda el: el == col, column_indexes))
        minors[idx] = det_serial(matrix, begin_row_index + 1, minor_cols)  # calculate required minors

    for j in range(n):
        result += sgn * matrix[begin_row_index][column_indexes[j]] * minors[j]
        sgn *= -1

    return result

if __name__ == "__main__":
    test_matrix = read_matrix("../../test_data/matrica5x5.txt")
    determinant = det_serial(test_matrix)
    print('det(matrix) =', determinant)





