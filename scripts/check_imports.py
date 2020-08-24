import subprocess
import sys


def check_imports():
    """
    Run `goimports -l .` to find all Go files whose imports aren't
    formatted appropriately.

    Imported Go packages don't only have to be ordered alphabetical,
    they also need to be grouped by imports from the standard library
    and from third-party vendors. For example:

    import (
        "errors"
        "os"

        "github.com/verless/verless/builder"
        "github.com/verless/verless/core/build"
        "github.com/verless/verless/parser"
    )

    If all files are formatted appropriately, exit with status 0,
    otherwise exit with status 1 by printing all invalid files.
    """
    output = subprocess.check_output(["goimports", "-l", "."])
    output = output.decode("utf-8")

    if output != "":
        msg = """goimports found non-ordered imports in the following files:
{0}
Run `goimports -w .` to order the imports or group and order them yourself.""".format(output)

        sys.exit(msg)
