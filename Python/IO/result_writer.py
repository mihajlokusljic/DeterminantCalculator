import os

class ExecutionResults:

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
    # Create results file if it does not exist, otherwise append to existing file
    with open(results_file_path, mode='a+') as res_file:
        # if file is empty, initialize the file with header
        if os.path.getsize(results_file_path) == 0:
            res_file.write("n,exec_time_ms,serial,implementation\n")

        # write new results
        for exec_res in execution_results:
                        res_file.write("{},{},{},{}\n"
                           .format(exec_res.n, exec_res.exec_time_ms, bool_to_string(exec_res.serial), "python"))

