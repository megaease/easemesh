package rcfile

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"gopkg.in/yaml.v2"
)

type (
	// RCFile contains information of rc file of emctl.
	RCFile struct {
		Server string `yaml:"server"`

		path string
	}
)

const (
	rcfileName = ".emctlrc"
)

// New creates an RCFile.
func New() (*RCFile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "get user home dir failed")
	}

	path := path.Join(homeDir, rcfileName)

	return &RCFile{
		path: path,
	}, nil
}

// Path returns the path of rc file.
func (r *RCFile) Path() string {
	return r.path
}

// Marshal marshals the content into rc file.
func (r *RCFile) Marshal() error {
	buff, err := yaml.Marshal(r)
	if err != nil {
		return errors.Wrapf(err, "marshal %+v to yaml failed", r)
	}

	err = ioutil.WriteFile(r.path, buff, 0644)
	if err != nil {
		return errors.Wrapf(err, "write file %s failed", r.path)
	}

	return nil
}

// Unmarshal Unmarshals the content from rc file.
func (r *RCFile) Unmarshal() error {
	buff, err := ioutil.ReadFile(r.path)
	if err != nil {
		return errors.Wrapf(err, "read file %s failed", r.path)
	}

	err = yaml.Unmarshal(buff, r)
	if err != nil {
		return errors.Wrapf(err, "unmarshal %s to yaml failed", buff)
	}

	return nil
}
