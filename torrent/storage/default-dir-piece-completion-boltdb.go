// Bolt piece completion is available.
//go:build !noboltdb && !cgo && !wasm
// +build !noboltdb,!cgo,!wasm

package storage

func NewDefaultPieceCompletionForDir(dir string) (PieceCompletion, error) {
	return NewBoltPieceCompletion(dir)
}
