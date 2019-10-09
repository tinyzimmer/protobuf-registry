package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const defaultVersion = "0.0.1"
const defaultRevision = "master"

type PostProtoRequest struct {
	ID            string             `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	Body          string             `json:"body,omitempty"`
	Version       string             `json:"version,omitempty"`
	RemoteDepends []*ProtoDependency `json:"remoteDeps,omitempty"`
}

type ProtoDependency struct {
	URL      string   `json:"url,omitempty"`
	Revision string   `json:"revision,omitempty"`
	Path     string   `json:"path,omitempty"`
	Ignores  []string `json:"ignores,omitempty"`
}

func NewProtoReqFromReader(rdr io.ReadCloser) (*PostProtoRequest, error) {
	defer rdr.Close()
	req := PostProtoRequest{
		Version:       defaultVersion,
		RemoteDepends: make([]*ProtoDependency, 0),
	}
	body, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &req)
	return &req, err
}

func (req *PostProtoRequest) Validate() error {
	if req.Name == "" || req.Body == "" {
		return errors.New("'name' and 'body' are required")
	}
	if len(req.RemoteDepends) > 0 {
		for _, depend := range req.RemoteDepends {
			if depend.URL == "" {
				return fmt.Errorf("Remote dependency URL cannot be blank")
			}
			if depend.Revision == "" {
				depend.Revision = defaultRevision
			}
		}
	}
	return nil
}
