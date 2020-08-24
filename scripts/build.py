import os
import shutil
import subprocess


def main():
    """Define a matrix of operating systems and architectures and
    build the verless binary for all permutations.

    If the current architecture is listed for an operating system
    in the exclusions dictionary, skip the permutation.
    """
    operating_systems = ["linux", "darwin", "windows"]
    architectures = ["386", "amd64", "arm"]
    exclusions = {
        "darwin": ["386", "arm"]
    }

    git_data = get_git_data()

    for go_os in operating_systems:
        for go_arch in architectures:

            if go_os in exclusions and go_arch in exclusions[go_os]:
                continue

            build(go_os, go_arch, git_data)
            package(go_os, go_arch)


def build(go_os, go_arch, git_data):
    """Build the verless binary for the given operating system and
    the given platform.

    The binary will be stored in target/<os>-<arch>. The binary
    name will be verless for Linux and macOS and verless.exe for
    Windows platforms.
    """
    binary = "verless.exe" if go_os == "windows" else "verless"
    target = "target/{0}-{1}/{2}".format(go_os, go_arch, binary)

    env = os.environ.copy()
    env["GOOS"] = go_os
    env["GOARCH"] = go_arch

    ld_flags = "-X github.com/verless/verless/config.GitTag={0} -X github.com/verless/verless/config.GitCommit={1}" \
        .format(git_data["tag"], git_data["commit"])

    subprocess.Popen(
        ["go", "build", "-v", "-ldflags", ld_flags, "-o", target, "cmd/verless/main.go"],
        env=env
    ).wait()


def package(go_os, go_arch):
    """
    Package a built binary as a zip or tar archive. It expects the
    binary in target/<os>-<arch>, where the build function stores
    its binaries.
    """
    ext = "zip" if go_os == "windows" else "tar"
    dest = "target/verless-{0}-{1}".format(go_os, go_arch)
    src = "target/{0}-{1}/".format(go_os, go_arch)

    shutil.make_archive(dest, ext, src)


def get_git_data():
    """
    Read the latest annotated Git tag without revision number as
    well as the short hash of the latest Git commit.

    The data is stored and returned as a dictionary.
    """
    git_tag = subprocess.check_output(["git", "describe", "--tags", "--abbrev=0"])
    git_tag = git_tag.decode("utf-8")

    git_commit = subprocess.check_output(["git", "rev-parse", "--short", "HEAD"])
    git_commit = git_commit.decode("utf-8")

    return {
        "tag": git_tag,
        "commit": git_commit,
    }


main()
