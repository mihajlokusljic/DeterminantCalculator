from IO.matrix_reader import read_matrix
from determinanat_calc.parallel_det_calc import det_parallel


MATRIX_PATH_TEMPLATE = "../../test_data/matrica{}x{}.txt"
RESULTS_BASE_PATH = "../../results"
STRONG_SCALING_MATRIX_ORDER = 10
AVAILABLE_MATRIX_ORDERS = [3, 5, 8, 9, 10, 11]

class Scaling_result:
    """
    Holds information about parallel determinant calculation to be stored in a results file.

    Attributes:
        matrix_order (int): order of matrix for which the determinant is calculated
        tasks_num (int): number of processes used during calculation
        exec_time_ms (float): duration of the calculation in milliseconds
    """
    def __init__(self, matrix_order, tasks_num, exec_time_ms):
        self.matrix_order = matrix_order
        self.tasks_num = tasks_num
        self.exec_time_ms = exec_time_ms

# Utility functions to load matrix and save execution results

def write_scaling_results(results_file_path, scaling_results):
    """
    Populates a given CSV file with given results. If file does not exists it will be created.
    The header contains: matrix_order (order of matrix for which the determinant is calculated),
    number_of_tasks (number of processes used in calculation), exec_time_ms (duration of the calculation
    in milliseconds).

    Args:
        results_file_path (string): path of the CSV file to hold results
        scaling_results (list(Scaling_result)): information about different determinant calculation attempts

    Return:
        None
    """
    with open(results_file_path, 'w') as file:
        # write header
        file.write("matrix_order,number_of_tasks,exec_time_ms\n")

        for result in scaling_results:
            file.write("{},{},{}\n"
                       .format(result.matrix_order, result.tasks_num, result.exec_time_ms))

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

def strong_scaling_experiment():
    """
    Loads a predefined matrix for strong scaling and calculates the determinant
    using parallel implementation with number of tasks ranging from one to the
    order of the matrix. Statistics about calculations are printed onto the console
    and written into a predefined results file: strong_scaling_results_python.csv.

    Return:
        None
    """
    print("Strong scaling experiment:\n")
    matrix = load_test_matrix(STRONG_SCALING_MATRIX_ORDER)
    results = []
    for tasks_num in range(1, STRONG_SCALING_MATRIX_ORDER + 1):
        determinant, exec_time_ms = det_parallel(matrix, tasks_num)
        results.append(Scaling_result(STRONG_SCALING_MATRIX_ORDER, tasks_num, exec_time_ms))

        print("Parallel determinant calculation of matrix of order {} with {} tasks took {} ms."
              .format(len(matrix), tasks_num, exec_time_ms))

    write_scaling_results("{}/strong_scaling_results_python.csv".format(RESULTS_BASE_PATH), results)

def weak_scaling_experiment():
    """
    Loads matrices of different orders from predefined files and calculates
    their determinants using parallel implementation where the number of tasks
    equals to the order of matrix. Statistics about calculations are printed onto the console
    and written into a predefined results file: weak_scaling_results_python.csv.

    Return:
        None
    """
    print("\nWeak scaling experiment:\n")

    results = []
    for n in AVAILABLE_MATRIX_ORDERS:
        matrix = load_test_matrix(n)
        determinant, exec_time_ms = det_parallel(matrix, n)
        results.append(Scaling_result(n, n, exec_time_ms))

        print("Parallel determinant calculation of matrix of order {} with {} tasks took {} ms."
              .format(len(matrix), n, exec_time_ms))

    write_scaling_results("{}/weak_scaling_results_python.csv".format(RESULTS_BASE_PATH), results)


if __name__ == "__main__":
    strong_scaling_experiment()
    weak_scaling_experiment()
