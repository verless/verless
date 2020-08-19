import subprocess


def check_imports():
    output = subprocess.check_output(["goimports"])
    output = output.decode("utf-8")

    print(output)