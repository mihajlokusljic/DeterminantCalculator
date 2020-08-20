import re

def read_matrix(matrix_file_path):
    """
    Loads and returns a regular matrix from given text file. The file must follow the convention:
    first row contains a non-negative integer number (n) which specifies the order of the matrix.
    The next n rows contain exactly n real numbers, separated by one or more space characters.

    Args:
        matrix_file_path(string): path of the file containing the matrix

    Return:
        matrix loaded into memory and ready for use (list(list(float))
    """
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