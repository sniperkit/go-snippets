package ipfs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/datatogether/core"
	"github.com/ipfs/go-datastore"
	// "github.com/jbenet/go-base58"
	// "github.com/multiformats/go-multihash"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

// Should be set by implementers
var IpfsApiServerUrl = ""

// TODO - add a skipHashed arg that allows us to skip urls that already have been seen
func ArchiveUrl(store datastore.Datastore, ipfsApiUrl string, url *core.Url) (headerHash, bodyHash string, err error) {
	urlstr := url.Url
	// header, body, err := GetUrlBytes(urlstr)
	// if err != nil {
	// 	err = fmt.Errorf("Error getting url: %s: %s", urlstr, err.Error())
	// 	return
	// }

	body, _, err := url.Get(store)
	if err != nil {
		err = fmt.Errorf("Error fetching url '%s': %s", urlstr, err.Error())
		return
	}

	// TODO - finish once we're storing the header hash with the url
	// if skipHashed && url.Hash != "" {
	// 	// skip hashing if we have a valid base58-encided multihash.
	// 	// some of our old, non-ipfs hashes used non-base58 encoded hashes
	// 	// this'll have the effect of considering them stale dated, which
	// 	// is what we we want.
	// 	if _, err := multihash.Decode(base58.Decode(url.Hash)); err == nil {
	// 		return
	// 	}
	// }

	buf := &bytes.Buffer{}
	if err = url.WarcRequest().Write(buf); err != nil {
		return
	}

	header := buf.Bytes()

	headerHash, err = WriteToIpfs(ipfsApiUrl, filepath.Base(urlstr), header)
	if err != nil {
		err = fmt.Errorf("Error writing %s header to ipfs: %s", filepath.Base(urlstr), err.Error())
		return
	}

	bodyHash, err = WriteToIpfs(ipfsApiUrl, filepath.Base(urlstr), body)
	if err != nil {
		err = fmt.Errorf("Error writing %s body to ipfs: %s", filepath.Base(urlstr), err.Error())
		return
	}

	// set hash for collection
	url.Hash = bodyHash

	if err = url.Save(store); err != nil {
		return
	}

	return
}

func WriteToIpfs(ipfsurl, filename string, data []byte) (hash string, err error) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	var (
		f       io.Writer
		ipfsReq *http.Request
		ipfsRes *http.Response
	)

	f, err = w.CreateFormFile("path", filename)
	if err != nil {
		err = fmt.Errorf("error creating form file: %s", err.Error())
		return
	}

	// TODO - handle errors
	if _, err = f.Write(data); err != nil {
		err = fmt.Errorf("error creating form file data: %s", err.Error())
		return
	}

	if err = w.Close(); err != nil {
		err = fmt.Errorf("error closing form file: %s", err.Error())
		return
	}

	// add to IPFS
	ipfsReq, err = http.NewRequest("POST", fmt.Sprintf("%s/add", ipfsurl), body)
	if err != nil {
		err = fmt.Errorf("error creating request: %s", err.Error())
		return
	}
	ipfsReq.Header.Set("Content-Type", w.FormDataContentType())

	ipfsRes, err = http.DefaultClient.Do(ipfsReq)
	if err != nil {
		err = fmt.Errorf("error sending request: %s", err.Error())
		return
	}
	defer ipfsRes.Body.Close()

	if ipfsRes.StatusCode != http.StatusOK {
		err = fmt.Errorf("IPFS Server returned non-200 status code: %d", ipfsRes.StatusCode)
		return
	}

	reply := struct {
		Reply string
		Hash  string
	}{}
	if err = json.NewDecoder(ipfsRes.Body).Decode(&reply); err != nil {
		err = fmt.Errorf("error reading response body: %s", err.Error())
		return
	}

	return reply.Hash, nil
}

func ReadFile(ipfsUrl, hash string) (io.ReadCloser, error) {
	res, err := http.Get(fmt.Sprintf("%s/cat?arg=%s", ipfsUrl, hash))
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
