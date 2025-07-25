package archive

import (
	"pm/internal/config"
	"testing"
)

func TestCreateArchive(t *testing.T) {
	type args struct {
		archivePath string
		targets     []config.Target
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				archivePath: "./test",
				targets:     []config.Target{{Path: "./*.go", Exclude: "test"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateArchive(tt.args.archivePath, tt.args.targets); (err != nil) != tt.wantErr {
				t.Errorf("CreateArchive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
}

func TestExtractArchive(t *testing.T) {
	type args struct {
		archivePath string
		dest        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				archivePath: "./test",
				dest:        "./exctracted/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExtractArchive(tt.args.archivePath, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("ExtractArchive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
