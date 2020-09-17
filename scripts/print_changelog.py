import subprocess

CHANGELOG_FILE = "CHANGELOG.md"


def print_changelog():
    """
    Print the changelog for the latest release tag by reading the
    changelog file. If the release tag does not exist as a release
    number in the changelog, the output will be empty.
    """
    git_tag = subprocess.check_output(["git", "describe", "--tags", "--abbrev=0"])
    git_tag = git_tag[:len(git_tag) - 1].decode("utf-8")

    # We'll start capturing the file contents when the heading for
    # the particular release appears, e.g. `## [1.2.3]`.
    start = "## [{0}]".format(git_tag[1:])

    # The `## [Unreleased]` heading will be ignored.
    unreleased = "## [Unreleased]"

    # We'll stop capturing the file content when the next release
    # heading appears, e.g. `## [1.2.2]`.
    end = "## ["

    capturing = False
    output = ""

    with open(CHANGELOG_FILE) as changelog:
        lines = changelog.readlines()

        for line in lines:
            # Start capturing if the line contains our release heading.
            if start in line and unreleased not in line:
                capturing = True
                continue
            # Stop capturing if we've reached the end, i.e. the next heading.
            if end in line and capturing:
                break
            if capturing:
                output += line

    print(output)


print_changelog()
