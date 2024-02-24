package util

import (
	"errors"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/tgragnato/magnetico/persistence"
)

func TotalSize(files []persistence.File) (uint64, error) {
	var totalSize uint64
	if len(files) == 0 {
		return 0, errors.New("no files would be persisted")
	}

	for _, file := range files {
		if file.Size < 0 {
			return 0, errors.New("file size less than zero")
		}

		totalSize += uint64(file.Size)
	}
	return totalSize, nil
}

func ValidateInfo(info *metainfo.Info) error {
	if len(info.Pieces)%20 != 0 {
		return errors.New("pieces has invalid length")
	}
	if info.PieceLength == 0 {
		return errors.New("zero piece length")
	}
	if int((info.TotalLength()+info.PieceLength-1)/info.PieceLength) != info.NumPieces() {
		return errors.New("piece count and file lengths are at odds")
	}
	return nil
}
