package decoder

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

type torrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

// TorrentFile struct for holding top level object
type TorrentFile struct {
	Announce   string
	Info       torrentInfo
	InfoHash   [20]byte
	PiecesHash [][20]byte
}

func calculateTorrentInfoHash(i *torrentInfo) ([20]byte, error) {

	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil

}

// calculate pieces hash of each one
func splitPieceHashes(i *torrentInfo) ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("bad pieces recieved", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
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

	infoHash, err := calculateTorrentInfoHash(&t.Info)

	if err != nil {
		return TorrentFile{}, err
	}

	t.PiecesHash, err = splitPieceHashes(&t.Info)
	if err != nil {
		return TorrentFile{}, err
	}

	t.InfoHash = infoHash
	return t, nil
}
