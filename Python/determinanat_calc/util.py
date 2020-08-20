import time

def measure_exec_time(decorated_func):

    def decorator(*args, **kwargs):
        start_time = time.time()
        result = decorated_func(*args, **kwargs)
        end_time = time.time()
        exec_time_ms = (end_time - start_time) * 1000
        return result, exec_time_ms

    return decorator