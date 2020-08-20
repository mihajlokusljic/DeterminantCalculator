import os

class ExecutionResults:
    """
    Contains information about execution times for determinant calculations, which
    are stored in the results file.

    Attributes:
        n (int): order of the matrix for which the determinant was calculated
        determinant (float): value of the determinant of matrix
        exec_time_ms (int): duration of the calculation process in milliseconds
        serial (bool): if true the implementation is serial, otherwise implementation is parallel
    """

    def __init__(self, matrix_order, determinant, exec_time_ms, serial):
        self.n = matrix_order
        self.determinant = determinant
        self.exec_time_ms = exec_time_ms
        self.serial = serial

def bool_to_string(bool):
    if bool:
        return "true"
    return "false"

def write_results(execution_results, results_file_path):
    """
    Writes information about execution times of calculating determinants of
    different matrices, using serial and parallel implementation, to a specified csv file.
    The CSV file header contains: n (matrix order), exec_time_ms (execution time), serial (determines
    whether the implementation is serial or parallel), implementation (determines the programming language
    of used implementation). If the file does not exist it will be created and initialized,
    otherwise the new results will be added to the end of existing file.

    Args:
        execution_results (list(ExecutionResults)): information about determinant calculation processes
        for different matrices and implementations
        results_file_path (string): path of the csv file to hold the informations

    Return:
        None
    """
    # Create results file if it does not exist, otherwise append to existing file
    with open(results_file_path, mode='a+') as res_file:
        # if file is empty, initialize the file with header
        if os.path.getsize(results_file_path) == 0:
            res_file.write("n,exec_time_ms,serial,implementation\n")

        # write new results
        for exec_res in execution_results:
                        res_file.write("{},{},{},{}\n"
                           .format(exec_res.n, exec_res.exec_time_ms, bool_to_string(exec_res.serial), "python"))

