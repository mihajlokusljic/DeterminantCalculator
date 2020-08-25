import time

def measure_exec_time(decorated_func):
    """
    Decorates a function by measuring the execution time of the function
    and returning the execution time along side with the result returned by the function.

    Args:
        decorated_func (function): a function whose execution time should be measured

    Return:
        a new function that also returns execution time of the original function (function)
    """

    def decorator(*args, **kwargs):
        start_time = time.time()
        result = decorated_func(*args, **kwargs)
        end_time = time.time()
        exec_time_ms = (end_time - start_time) * 1000
        try:
            ret = []
            for res_value in result:
                ret.append(res_value)
            ret.append(exec_time_ms)
            return ret
        except:
            return [result, exec_time_ms]

    return decorator