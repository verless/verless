import shutil
import subprocess


def matrix():
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

    for os in operating_systems:
        for arch in architectures:

            if os in exclusions and arch in exclusions[os]:
                continue

            build(os, arch)
            package(os, arch)


def build(os, arch):
    """Build the verless binary for the given operating system and
    the given platform.

    The binary will be stored in ../target/<os>-<arch>. The binary
    name will be verless for Linux and macOS and verless.exe for
    Windows platforms.
    """
    go_os = "GOOS={0}".format(os)
    go_arch = "GOARCH={0}".format(arch)

    binary = "verless.exe" if os == "windows" else "verless"
    target = "../target/{0}-{1}/{2}".format(os, arch, binary)

    subprocess.run([go_os, go_arch, "go", "build", "-v", "-o", target, "../cmd/verless/main.go"])


def package(os, arch):
    ext = "zip" if os == "windows" else "tar"
    src = "../target/{0}-{1}".format(os, arch)
    dest = "../target/{0}-{1}".format(os, arch)

    shutil.make_archive(src, ext, dest)


matrix()
