package decoder

import (
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type torrentInfo struct {
	Name        string
	Length      int
	PieceLength int `bencode:"piece length"`
	Pieces      string
}

// TorrentFile struct for holding top level object
type TorrentFile struct {
	Announce string
	Info     torrentInfo
}

// Decode decode the file and return TorrentFile object
func Decode(f string) (TorrentFile, error) {

	data, err := os.Open(f)
	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{}
	err = bencode.Unmarshal(data, &t)
	if err != nil {
		return TorrentFile{}, err
	}
	return t, nil
}