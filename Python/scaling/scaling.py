from IO.matrix_reader import read_matrix
from determinanat_calc.parallel_det_calc import det_parallel
from determinanat_calc.serial_det_calc import det_serial


MATRIX_PATH_TEMPLATE = "../../test_data/matrica{}x{}.txt"
RESULTS_BASE_PATH = "../../results"
STRONG_SCALING_MATRIX_ORDER = 10
AVAILABLE_MATRIX_ORDERS = [3, 5, 8, 9, 10, 11]

SERIAL_CODE_SHARE = 0
PARALLEL_CODE_SHARE = 1

class Scaling_result:
    """
    Holds information about determinant calculation for matrix of certain order to be stored in a results file.

    Attributes:
        matrix_order (int): order of matrix for which the determinant is calculated
        serial_exec_time_ms (float): execution time of serial implementation in milliseconds
        parallel_tasks_num (int): number of processes used in parallel calculation)
        parallel_exec_time_ms (float): execution time of parallel implementation in milliseconds
        achieved_speedup (float): actual speedup calculated using serial and parallel execution time
        max_theoretical_speedup (float): maximum possible speedup considering serial and parallel portions of code
    """
    def __init__(self, matrix_order, serial_exec_time_ms, parallel_tasks_num, parallel_exec_time_ms,
                 achieved_speedup, max_speedup):
        self.matrix_order = matrix_order
        self.serial_exec_time_ms = serial_exec_time_ms
        self.parallel_tasks_num = parallel_tasks_num
        self.parallel_exec_time_ms = parallel_exec_time_ms
        self.achieved_speedup = achieved_speedup
        self.max_theoretical_speedup = max_speedup

# Utility functions to load matrix and save execution results

def write_scaling_results(results_file_path, scaling_results):
    """
    Populates a given CSV file with given results. If file does not exists it will be created.
    The header contains: matrix_order (order of matrix for which the determinant is calculated),
    serial_exec_time_ms (execution time of serial implementation in milliseconds),
    parallel_tasks_num (number of processes used in parallel calculation), parallel_exec_time_ms
    (execution time of parallel implementation in milliseconds), achieved_speedup (actual speedup calculated using
    serial and parallel execution time), max_theoretical_speedup (maximum possible speedup considering serial and
    parallel portions of code).

    Args:
        results_file_path (string): path of the CSV file to hold results
        scaling_results (list(Scaling_result)): information about different determinant calculation attempts

    Return:
        None
    """
    with open(results_file_path, 'w') as file:
        # write header
        file.write("matrix_order,serial_exec_time_ms,parallel_tasks_num,parallel_exec_time_ms,achieved_speedup,max_theoretical_speedup\n")

        for result in scaling_results:
            file.write("{},{},{},{},{},{}\n"
                       .format(result.matrix_order, result.serial_exec_time_ms, result.parallel_tasks_num,
                               result.parallel_exec_time_ms, result.achieved_speedup, result.max_theoretical_speedup))

def load_test_matrix(matrix_order):
    """
    Loads a test matrix of given order from a predefined directory.

    Args:
        matrix_order (int): order of desired matrix

    Return:
        matrix loaded into memory (list(list(float)))
    """
    matrix_path = MATRIX_PATH_TEMPLATE.format(matrix_order, matrix_order)
    matrix = read_matrix(matrix_path)
    print( "Successfully loaded matrix of order {}. Calculating determinant(s)...".format(len(matrix)) )
    return matrix

# Scaling experiment functions

def max_speedup_Amdahl(processes_num):
    return 1.0 / (SERIAL_CODE_SHARE + PARALLEL_CODE_SHARE / processes_num)

def max_speedup_Gustafson(processes_num):
    return SERIAL_CODE_SHARE + PARALLEL_CODE_SHARE * processes_num

def strong_scaling_experiment():
    """
    Loads a predefined matrix for strong scaling and calculates the determinant
    using parallel implementation with number of tasks ranging from one to the
    order of the matrix. Statistics about calculations are printed onto the console
    and written into a predefined results file: strong_scaling_results_python.csv.

    Return:
        None
    """
    print("==========================")
    print("Strong scaling experiment:")
    print("==========================\n\n")
    matrix = load_test_matrix(STRONG_SCALING_MATRIX_ORDER)

    determinant, potential_parallel_code_exec_time_ms, serial_exec_time_ms = det_serial(matrix)
    PARALLEL_CODE_SHARE = potential_parallel_code_exec_time_ms / serial_exec_time_ms
    SERIAL_CODE_SHARE = 1 - PARALLEL_CODE_SHARE
    print("Serial determinant calculation of matrix of order {} took {} ms."
          .format(len(matrix), serial_exec_time_ms))
    print("Execution of code that can be run in parallel took {} ms.".format(potential_parallel_code_exec_time_ms))
    print("Parallel code share is: {}\nSerial code share is: {}\n".format(PARALLEL_CODE_SHARE, SERIAL_CODE_SHARE))


    results = []
    for tasks_num in range(2, STRONG_SCALING_MATRIX_ORDER + 1):
        determinant, parallel_exec_time_ms = det_parallel(matrix, tasks_num)
        achieved_speedup = serial_exec_time_ms / parallel_exec_time_ms
        max_speedup = max_speedup_Amdahl(tasks_num)
        result = Scaling_result(STRONG_SCALING_MATRIX_ORDER, serial_exec_time_ms, tasks_num,
                                parallel_exec_time_ms, achieved_speedup, max_speedup)
        results.append(result)

        print("Parallel determinant calculation of matrix of order {} with {} tasks took {} ms."
              .format(STRONG_SCALING_MATRIX_ORDER, tasks_num, parallel_exec_time_ms))
        print("Achieved speedup is: {}X.\nMaximum speedup according to Amdahl’s law is: {}X.\n".format(achieved_speedup, max_speedup))

    write_scaling_results("{}/strong_scaling_results_python.csv".format(RESULTS_BASE_PATH), results)
    print("Successfully finished strong scaling experiment.")

def weak_scaling_experiment():
    """
    Loads matrices of different orders from predefined files and calculates
    their determinants using parallel implementation where the number of tasks is
    equal to the order of matrix. Statistics about calculations are printed onto the console
    and written into a predefined results file: weak_scaling_results_python.csv.

    Return:
        None
    """
    print("\n\n========================")
    print("Weak scaling experiment:")
    print("========================\n\n")


    results = []
    for n in AVAILABLE_MATRIX_ORDERS:
        matrix = load_test_matrix(n)
        determinant, _, serial_exec_time_ms = det_serial(matrix)
        print("Serial determinant calculation of matrix of order {} took {} ms."
              .format(len(matrix), serial_exec_time_ms))

        determinant, parallel_exec_time_ms = det_parallel(matrix, n)
        achieved_speedup = serial_exec_time_ms / parallel_exec_time_ms
        max_speedup = max_speedup_Gustafson(n)
        result = Scaling_result(n, serial_exec_time_ms, n, parallel_exec_time_ms, achieved_speedup, max_speedup)
        results.append(result)
        print("Parallel determinant calculation of matrix of order {} with {} tasks took {} ms."
              .format(n, n, parallel_exec_time_ms))
        print("Achieved speedup is: {}X.\nMaximum speedup according to Gustafson’s law is: {}X.\n".format(achieved_speedup, max_speedup))

    write_scaling_results("{}/weak_scaling_results_python.csv".format(RESULTS_BASE_PATH), results)
    print("Successfully finished weak scaling experiment.")


if __name__ == "__main__":
    strong_scaling_experiment()
    weak_scaling_experiment()
