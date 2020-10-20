package fs

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/verless/verless/test"
)

func TestIsSafeToRemove(t *testing.T) {
	// Create new in-memory FS with one file and one directory
	tempFS := afero.NewMemMapFs()
	_, _ = tempFS.Create("file.go")
	_ = tempFS.Mkdir("directory", 0755)

	type args struct {
		targetFs afero.Fs
		path     string
		force    bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Add test entries here
		{"Remove file", args{tempFS, "file.go", false}, false},
		{"Remove directory", args{tempFS, "directory", false}, false},
		{"Remove non-existent file", args{tempFS, "notexist", false}, true},
		{"Force remove", args{tempFS, "file.go", true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSafeToRemove(tt.args.targetFs, tt.args.path, tt.args.force)
			test.Equals(t, got, tt.want)
		})
	}
}

func TestRmdir(t *testing.T) {
	// Create new in-memory FS with -
	// One file, One empty directory, One directory with files
	tempFS := afero.NewMemMapFs()
	_, _ = tempFS.Create("file.go")
	_ = tempFS.Mkdir("emptydir", 0755)
	_ = tempFS.MkdirAll("first/second", 0755)
	_, _ = tempFS.Create("first/second/third.go")

	type args struct {
		fs   afero.Fs
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test entries here
		{"Remove empty directory", args{tempFS, "emptydir"}, false},
		{"Remove file", args{tempFS, "file.go"}, false},
		{"Remove directory recursively", args{tempFS, "first"}, false},
		{"Try on non-existent file", args{tempFS, "notexist.go"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Rmdir(tt.args.fs, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Rmdir() error = %v, wantErr %v", err, tt.wantErr)
				test.ExpectedError(t, nil, err)
			}
		})
	}
}

func TestCopyFromOS(t *testing.T) {
	// Create blank in-memory FS
	// One directory for source files, one for destination files
	// One nested directory for source files being copied with structure
	// One file in the root directory to copy into the destination
	// One empty directory, one with no read permissions
	tempFS := afero.NewMemMapFs()
	_ = tempFS.Mkdir("src", 0755)
	_, _ = tempFS.Create("src/srcfile1.go")
	_ = tempFS.Mkdir("dest", 0755)
	_ = tempFS.MkdirAll("nestdir/levelone", 0755)
	_, _ = tempFS.Create("nestdir/levelone/srcfile2.go")
	_, _ = tempFS.Create("srcfile3.go")
	_ = tempFS.Mkdir("empty", 0755)
	_ = tempFS.Mkdir("restricted", 0000)
	_, _ = tempFS.Create("restricted/cantread.go")

	type args struct {
		targetFs afero.Fs
		src      string
		dest     string
		fileOnly bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test entries here
		{"Copy from src to dest dir", args{tempFS, "src", "dest", true}, false},
		{"Copy with source structure", args{tempFS, "nestdir", "dest", false}, false},
		{"Copy file from parent of dest", args{tempFS, "srcfile3.go", "dest", false}, false},
		{"Destination directory does not exist", args{tempFS, "nestdir", "newdest", false}, false},
		{"Source directory is empty", args{tempFS, "empty", "dest", false}, false},
		{"Source file does not exist", args{tempFS, "nofile.go", "dest", true}, false},
		{"Source directory has no read perms", args{tempFS, "restricted", "dest", false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFromOS(tt.args.targetFs, tt.args.src, tt.args.dest, tt.args.fileOnly); (err != nil) != tt.wantErr {
				t.Errorf("CopyFromOS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamFiles(t *testing.T) {
	// Filter functions kept here for independence from fs.go

	// MarkdownOnly is a filter that only lets pass Markdown files.
	MarkdownOnly = func(file string) bool {
		return filepath.Ext(file) == ".md"
	}

	// NoUnderscores is a predefined filter that doesn't let pass
	// files starting with an underscore.
	NoUnderscores = func(file string) bool {
		filename := filepath.Base(file)
		return !strings.HasPrefix(filename, "_")
	}

	// Filters:
	noFilter := []func(string) bool{}
	markFilter := []func(string) bool{MarkdownOnly}
	underscoreFilter := []func(string) bool{NoUnderscores}
	allFilter := []func(string) bool{MarkdownOnly, NoUnderscores}

	type args struct {
		path    string
		files   chan string
		filters []func(file string) bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test entries here
		{"No filters", args{"../example", make(chan string), noFilter}, false},
		{"Markdown filters", args{"../example", make(chan string), markFilter}, false},
		{"Underscore filters", args{"../example", make(chan string), underscoreFilter}, false},
		{"All filters", args{"../example", make(chan string), allFilter}, false},
		{"Path does not exist", args{"notexist", make(chan string), allFilter}, false},
	}
	for _, tt := range tests {
		go func() {
			// Drain all data that goes into channel
			for file := range tt.args.files {
				_ = file
			}
		}()
		err := StreamFiles(tt.args.path, tt.args.files, tt.args.filters...)
		if (err != nil) != tt.wantErr {
			t.Errorf("StreamFiles() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}
