package fs

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/verless/verless/test"
)

func TestIsSafeToRemove(t *testing.T) {
	type args struct {
		targetFs afero.Fs
		path     string
		force    bool
	}

	tempFS := afero.NewMemMapFs()
	tests := []struct {
		name string
		args args
		want bool
	}{
		// Add test entries here
		{"Remove Example", args{tempFS, "../example", false}, true},
		{"Should not remove", args{tempFS, ".", false}, false},
		{"Force Remove", args{tempFS, ".", true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSafeToRemove(tt.args.targetFs, tt.args.path, tt.args.force)
			test.Equals(t, got, tt.want)
		})
	}
}

func TestRmdir(t *testing.T) {
	type args struct {
		fs   afero.Fs
		path string
	}
	tempFS := afero.NewMemMapFs()
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// Add test entries here
		{"Remove example", args{tempFS, "../example"}, false},
		// BUG: This works, but it should not, according to IsSafeToRemove
		{"Should not remove", args{tempFS, "."}, false},
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
