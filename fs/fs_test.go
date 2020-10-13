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
