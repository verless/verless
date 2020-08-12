import subprocess


def matrix():
    operating_systems = ["linux", "darwin", "windows"]
    architectures = ["386", "amd64", "arm"]

    for os in operating_systems:
        for arch in architectures:
            build(os, arch)


def build(os, arch):
    go_os = "GOOS={0}".format(os)
    go_arch = "GOARCH={0}".format(arch)
    target = "../target/{0}-{1}/verless".format(os, arch)
    subprocess.run([go_os, go_arch, "go", "build", "-v", "-o", target, "../cmd/main.go"])


matrix()
