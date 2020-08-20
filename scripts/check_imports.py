import subprocess
import sys


def check_imports():
    output = subprocess.check_output(["goimports", "-l", ".."])
    output = output.decode("utf-8")

    if output != "":
        msg = """goimports found non-ordered imports in the following files:
{0}
Run `goimports -w .` to order the imports or group and order them yourself.""".format(output)

        sys.exit(msg)
