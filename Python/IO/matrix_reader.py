import re

def read_matrix(matrix_file_path):
    matrix = []

    with open(matrix_file_path, 'r') as in_file:
        n = int(in_file.readline())
        for row_idx in range(n):
            row = in_file.readline().strip()
            # columns are separated with single or multiple space characters
            row_cols = [float(num) for num in re.split(' +', row)]
            if len(row_cols) != n:
                raise Exception("Sqared matrix is required. Found {} columns in row {}, expected {}."
                                .format(len(row_cols), row_idx + 1, n))
            matrix.append(row_cols)

    return matrix

if __name__ == "__main__":
    test_matrix = read_matrix("../../Test data/matrica9x9.txt")
    assert len(test_matrix) == 9
    assert test_matrix[3][2] == -24.9