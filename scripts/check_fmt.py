import subprocess
import sys


def check_fmt():
    """Run `gofmt -l ..` to find all Go files that aren't formatted
    appropriately.

    If all files are formatted appropriately, exit with status 0,
    otherwise exit with status 1 by printing all invalid files.
    """
    output = subprocess.check_output(["gofmt", "-l", ".."])
    output = output.decode("utf-8")

    if output != "":
        msg = """gofmt found non-formatted code in the following files:
{0}
Run `go fmt ./...` to format your files and commit them again.""".format(output)

        sys.exit(msg)


check_fmt()
