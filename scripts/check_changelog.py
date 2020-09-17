import subprocess
import sys

CHANGELOG_FILE = "CHANGELOG.md"


def check_changelog():
    """
    When a new release tag is pushed in order to create a new verless
    release, check if that particular release is mentioned in the
    changelog.

    For a release tag like `v1.2.3`, the changelog has to contain a
    release section called `[1.2.3]`. If the release isn't mentioned
    in the changelog, exit with an error.
    """
    git_tag = subprocess.check_output(["git", "describe", "--tags", "--abbrev=0"])

    # The output from subprocess.check_output() is b'v1.2.3\n'. First,
    # cut off the new line \n and decode it to UTF-8.
    git_tag = git_tag[:len(git_tag) - 1].decode("utf-8")

    # Cut off the `v` prefix to get the actual release number.
    search_expr = "[{0}]".format(git_tag[1:])

    with open(CHANGELOG_FILE) as changelog:
        content = changelog.read()
        if search_expr not in content:
            msg = """You're trying to create a new release tag {0}, but that release is not mentioned
in the changelog. Add a section called {1} to {2} and push again.""" \
                .format(git_tag, search_expr, CHANGELOG_FILE)

            sys.exit(msg)


check_changelog()
