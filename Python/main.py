import sys
import time

from determinanat_calc.serial_det_calc import det_serial
from determinanat_calc.parallel_det_calc import det_parallel
from IO.matrix_reader import read_matrix
from IO.result_writer import ExecutionResults, write_results

def execute_serial_calculation(matrix, matrix_file_path):
    start_time = time.time()
    determinant = det_serial(matrix)
    end_time = time.time()
    exec_time_ms = (end_time - start_time) * 1000

    return ExecutionResults(matrix_file_path, len(matrix), determinant, exec_time_ms, True)

def execute_parallel_calculation(matrix, matrix_file_path):
    start_time = time.time()
    determinant = det_parallel(matrix)
    end_time = time.time()
    exec_time_ms = (end_time - start_time) * 1000

    return ExecutionResults(matrix_file_path, len(matrix), determinant, exec_time_ms, False)


if __name__ == "__main__":
    arguments = sys.argv[1:]
    results_file = arguments[-1]
    execution_results = []
    for matrix_file_path in arguments[:-1]:
        matrix = read_matrix(matrix_file_path)
        print("\nSuccessfully loaded matrix: {}\nCalculating determinants...\n".format(matrix_file_path))

        # Serial calculation
        serial_result = execute_serial_calculation(matrix, matrix_file_path)
        execution_results.append(serial_result)
        print("Serial calculation result:\ndet(mat) = {}\nSerial calculation took {} ms.\n"
              .format(serial_result.determinant, serial_result.exec_time_ms))

        # Parallel calculation
        parallel_result = execute_parallel_calculation(matrix, matrix_file_path)
        execution_results.append(parallel_result)
        print("Parallel calculation result:\ndet(mat) = {}\nParallel calculation took {} ms.\n"
              .format(parallel_result.determinant, parallel_result.exec_time_ms))


        if len(arguments) > 2:
            print("=" * 20)

    write_results(execution_results, arguments[-1])

