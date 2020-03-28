package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	bencode "github.com/jackpal/bencode-go"
	"github.com/khubo/go-torrent/decoder"
	peer "github.com/khubo/go-torrent/peer"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func constructTrackerURL(t decoder.TorrentFile, peerID [20]byte, port int) (string, error) {
	base, err := url.Parse(t.Announce)

	if err != nil {
		fmt.Println("Errored")
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

// RequestPeers method to request for torrent peers
func RequestPeers(t decoder.TorrentFile, id [20]byte, port int) ([]peer.Peer, error) {
	url, err := constructTrackerURL(t, id, port)

	if err != nil {
		log.Fatalln(err)
	}
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		return nil, err
	}

	peers, err := peer.ExtractPeers([]byte(trackerResp.Peers))
	if err != nil {
		return nil, err
	}

	return peers, nil
}