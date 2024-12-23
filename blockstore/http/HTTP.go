package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/metux/go-nebulon/blockstore/common"
	"github.com/metux/go-nebulon/core/base"
	"github.com/metux/go-nebulon/core/wire"
)

type HTTP struct {
	common.StoreBase
	Url   string
	Error error
}

func (s HTTP) PutBlock(data []byte, reftype wire.RefType) (base.BlockRef, error) {
	ref := wire.RefForBlock(data, reftype)
	code, _, hdr, err := s.blockRefUrlCall(http.MethodPut, ref, data, nil)

	if err != nil {
		return ref, err
	}

	newref := wire.ParseBlockRef(hdr.Get(wire.HttpV1_Header_BlockRef))
	// FIXME: compare newref and oldref

	switch code {
	case http.StatusOK, http.StatusCreated:
		return newref, nil
	default:
		return newref, wire.HttpV1_Error(code)
	}
}

func (s HTTP) request1(method string, url string, data []byte, hdr http.Header) (*http.Response, error) {
	client := http.Client{}

	req, err := http.NewRequest(method, s.Url+"/"+url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header = hdr

	return client.Do(req)
}

func (s HTTP) request(method string, url string, data []byte, hdr http.Header) (int, []byte, http.Header, error) {
	resp, err := s.request1(method, url, data, hdr)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, body, resp.Header, err
}

func refUrlPart(ref base.BlockRef) string {
	return ref.Type.String() + "/" + ref.OID()
}

func (s HTTP) blockRefUrlCall(method string, ref base.BlockRef, data []byte, hdr http.Header) (int, []byte, http.Header, error) {
	return s.request(method, "v1/block/"+refUrlPart(ref), data, hdr)
}

func (s HTTP) blockInfUrlCall(method string, ref base.BlockRef, data []byte, hdr http.Header) (int, []byte, http.Header, error) {
	return s.request(method, "v1/blockinf/"+refUrlPart(ref), data, hdr)
}

func (s HTTP) GetBlock(ref base.BlockRef) ([]byte, error) {
	status, body, _, err := s.blockRefUrlCall(http.MethodGet, ref, nil, nil)
	if err != nil {
		return nil, err
	}
	switch status {
	case http.StatusOK:
		return body, nil
	case http.StatusNotFound:
		return nil, base.ErrNotFound
	default:
		return body, wire.HttpV1_Error(status)
	}
}

func (s HTTP) KeepBlock(ref base.BlockRef) error {
	status, _, _, err := s.blockRefUrlCall(http.MethodPatch, ref, nil, nil)

	if err != nil {
		return err
	}

	switch status {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return base.ErrNotFound
	default:
		return wire.HttpV1_Error(status)
	}

	return err
}

func (s HTTP) IterateBlocks() base.BlockRefStream {
	ch := make(base.BlockRefStream, IterateChanSize)

	go func() {
		resp, err := s.request1(http.MethodGet, "v1/block", nil, http.Header{})
		if err != nil {
			ch.SendError(err)
			close(ch)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			t := scanner.Text()
			ref := wire.ParseBlockRef(t)
			ch.SendRef(ref)
		}
		if err := scanner.Err(); err != nil {
			ch.SendError(err)
		}
		ch.Finish()
		resp.Body.Close()
	}()

	return ch
}

func (s HTTP) DeleteBlock(ref base.BlockRef) error {
	status, _, _, err := s.blockRefUrlCall(http.MethodDelete, ref, nil, nil)

	if err != nil {
		return err
	}

	switch status {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusNotFound:
		return base.ErrNotFound
	default:
		return wire.HttpV1_Error(status)
	}
}

func (s HTTP) PeekBlock(ref base.BlockRef, fetch int) (base.BlockInfo, error) {
	inf := base.BlockInfo{Ref: ref}

	code, body, _, err := s.request(http.MethodGet, "v1/blockinf/"+refUrlPart(ref), nil,
		http.Header{wire.HttpV1_Header_FetchDepth: []string{strconv.Itoa(fetch)}})

	if err = json.Unmarshal(body, &inf); err != nil {
		return inf, err
	}

	switch code {
	case http.StatusNotFound:
		return inf, base.ErrNotFound
	case http.StatusOK:
		return inf, nil
	default:
		return inf, wire.HttpV1_Error(code)
	}
}

func (s HTTP) Ping() error {
	return nil
}

func NewByConfig(config base.BlockStoreConfig, links map[string]base.IBlockStore) (*HTTP, error) {
	if config.Url == "" {
		return nil, base.ErrMissingUrl
	}
	return &HTTP{
		StoreBase: common.StoreBase{
			Name: config.ID(),
		},
		Url: config.Url,
	}, nil
}
