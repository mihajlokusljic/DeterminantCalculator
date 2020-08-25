import sys

from determinanat_calc.serial_det_calc import det_serial
from determinanat_calc.parallel_det_calc import det_parallel
from IO.matrix_reader import read_matrix
from IO.result_writer import ExecutionResults, write_results

def execute_serial_calculation(matrix):
    """
    Calculates the determinant of given regular matrix using serial implementation
    and measures the execution time.

    Args:
        matrix(list(list(float)): regular matrix for which the determinant is calculated

    Return:
        result object that can be stored in the results file (ExecutionResults)
    """
    determinant, _, exec_time_ms = det_serial(matrix)

    return ExecutionResults(len(matrix), determinant, exec_time_ms, True)

def execute_parallel_calculation(matrix):
    """
    Calculates the determinant od given regular matrix using parallel implementation
    and measures the execution time.

    Args:
        matrix(list(list(float)): regular matrix for which the determinant is calculated

    Return:
        result object that can be stored in the results file (ExecutionResults)
    """
    determinant, exec_time_ms = det_parallel(matrix, len(matrix))

    return ExecutionResults(len(matrix), determinant, exec_time_ms, False)

def process_matrices():
    """
    This is the main method. It will load matrices from files specified in arguments to the program,
    calculate the determinant of each matrix, both sequentially and parallel, and append data about
    the execution times to a results file, specified as the last argument to the program. Results
    and execution times are also printed to the console.
    The first expected argument is the program itself. The next n arguments should be paths to
    valid text files containing a regular matrix: first row contains a number (n) which specifies the
    order of the matrix. The next n rows contain exactly n real numbers, separated by one or more
    space characters. The last argument to the program is a path to a csv file where information
    about the execution times for different matrices are stored.

    Return:
        None
    """

    arguments = sys.argv[1:]
    results_file = arguments[-1]
    execution_results = []

    for matrix_file_path in arguments[:-1]:
        # load matrix
        matrix = read_matrix(matrix_file_path)
        print("\nSuccessfully loaded matrix: {}\nCalculating determinants...\n".format(matrix_file_path))

        # Serial calculation
        serial_result = execute_serial_calculation(matrix)
        execution_results.append(serial_result)
        print("Serial calculation result:\ndet(mat) = {}\nSerial calculation took {} ms.\n"
              .format(serial_result.determinant, serial_result.exec_time_ms))

        # Parallel calculation
        parallel_result = execute_parallel_calculation(matrix)
        execution_results.append(parallel_result)
        print("Parallel calculation result:\ndet(mat) = {}\nParallel calculation took {} ms.\n"
              .format(parallel_result.determinant, parallel_result.exec_time_ms))

        # Separate console output
        if len(arguments) > 2:
            print("=" * 20)

    write_results(execution_results, results_file)

if __name__ == "__main__":
    process_matrices()