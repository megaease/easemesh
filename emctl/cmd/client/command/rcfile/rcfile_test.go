package rcfile

import (
	"os"
	"path"
	"testing"

	utiltesting "k8s.io/client-go/util/testing"
)

func TestRCFile(t *testing.T) {
	tmpDir, err := utiltesting.MkTmpdir("rcfile")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	expectRCFile := RCFile{path: path.Join(tmpDir, rcfileName), Server: "127.0.0.1:3333"}

	err = expectRCFile.Marshal()
	if err != nil {
		t.Fatalf("marshal %+v failed %s", expectRCFile, err)
	}

	_, err = os.Stat(path.Join(tmpDir, rcfileName))
	if err != nil {
		os.Remove(expectRCFile.path)
		t.Fatalf("marshal emctlrc file error: %s", err)
	}

	rcFile := RCFile{path: expectRCFile.path}

	err = rcFile.Unmarshal()
	if err != nil {
		os.Remove(expectRCFile.path)
		t.Fatalf("unmarshal rc %+v file error: %s", rcFile, err)
	}

	if rcFile.Server != expectRCFile.Server {
		os.Remove(expectRCFile.path)
		t.Fatalf("expect server: [%s] but: [%s]", expectRCFile.Server, rcFile.Server)
	}

	os.Remove(expectRCFile.path)
}

func TestRcFileMarshalShouldError(t *testing.T) {
	r := RCFile{}
	err := r.Marshal()
	if err == nil {
		t.Fatalf("marsh r %+v should error", r)
	}
}

func TestRCNew(t *testing.T) {
	rc, err := New()
	if err != nil {
		t.Fatalf("new rcfile error: %s", err)
	}

	homeDir, _ := os.UserHomeDir()
	expectPath := path.Join(homeDir, rcfileName)
	if rc.Path() != expectPath {
		t.Fatalf("expect rc path %s but %s", expectPath, rc.path)
	}
}
