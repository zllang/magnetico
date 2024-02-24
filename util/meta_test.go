package util_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/tgragnato/magnetico/persistence"
	"github.com/tgragnato/magnetico/util"
)

func TestTotalSize(t *testing.T) {
	t.Parallel()
	positiveRand := rand.Int63n(math.MaxInt64)

	tests := []struct {
		name    string
		files   []persistence.File
		want    uint64
		wantErr bool
	}{
		{
			name:    "No elements",
			files:   []persistence.File{},
			want:    0,
			wantErr: true,
		},
		{
			name: "Negative size",
			files: []persistence.File{
				{
					Size: -rand.Int63n(math.MaxInt64),
					Path: "",
				},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Zero size",
			files: []persistence.File{
				{
					Size: 0,
					Path: "",
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Positive size",
			files: []persistence.File{
				{
					Size: positiveRand,
					Path: "",
				},
			},
			want:    uint64(positiveRand),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.TotalSize(tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("TotalSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TotalSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		info    *metainfo.Info
		wantErr bool
	}{
		{
			name:    "Empty test",
			info:    &metainfo.Info{},
			wantErr: true,
		},
		{
			name: "Pieces not %20",
			info: &metainfo.Info{
				PieceLength: 5,
				Pieces:      []byte{0, 0, 0, 0, 0},
				Name:        "",
				NameUtf8:    "",
				Length:      0,
				Private:     nil,
				Source:      "",
				Files:       []metainfo.FileInfo{},
			},
			wantErr: true,
		},
		{
			name: "Invalid file length",
			info: &metainfo.Info{
				PieceLength: 1,
				Pieces:      []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Name:        "",
				NameUtf8:    "",
				Length:      0,
				Private:     nil,
				Source:      "",
				Files:       []metainfo.FileInfo{},
			},
			wantErr: true,
		},
		{
			name: "valid info",
			info: &metainfo.Info{
				Pieces:      make([]byte, 20),
				PieceLength: 1,
				Length:      20,
				Files:       []metainfo.FileInfo{{Length: 1, Path: []string{"file1"}}},
			},
			wantErr: false,
		},
		{
			name: "invalid pieces length",
			info: &metainfo.Info{
				Pieces:      make([]byte, 21),
				PieceLength: 1,
				Length:      20,
			},
			wantErr: true,
		},
		{
			name: "zero piece length with total length",
			info: &metainfo.Info{
				Pieces:      make([]byte, 20),
				PieceLength: 0,
				Length:      20,
			},
			wantErr: true,
		},
		{
			name: "mismatch piece count and file lengths",
			info: &metainfo.Info{
				Pieces:      make([]byte, 20),
				PieceLength: 1,
				Length:      21,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := util.ValidateInfo(tt.info); (err != nil) != tt.wantErr {
				t.Errorf("ValidateInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
