//    Copyright 2016 Red Hat, Inc.
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package download_test

import (
	"bytes"
	"crypto"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	download "github.com/jimmidyson/go-download"
)

func TestDownloadToFileFailOnMkdirs(t *testing.T) {
	err := download.ToFile("http://whatever:12345", "non-existent-directory", download.FileOptions{Mkdirs: download.MkdirNone})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDownloadToFileSuccess(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testData, err := ioutil.ReadFile(filepath.Join("testdata", "testfile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	downloadedData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(testData, downloadedData) {
		t.Fatal("wrong downloaded data")
	}
}

func TestDownloadToFileSuccessMkdirs(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpDir, err := ioutil.TempDir(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpDir) }()
	_ = os.Remove(tmpDir)

	tmpFile := filepath.Join(tmpDir, "tmp")
	err = download.ToFile(srv.URL+"/testfile", tmpFile, download.FileOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testData, err := ioutil.ReadFile(filepath.Join("testdata", "testfile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	downloadedData, err := ioutil.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(testData, downloadedData) {
		t.Fatal("wrong downloaded data")
	}
}

func TestDownloadToFileSuccessMD5Checksum(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{
		Options: download.Options{
			Checksum:     "d577273ff885c3f84dadb8578bb41399",
			ChecksumHash: crypto.MD5,
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testData, err := ioutil.ReadFile(filepath.Join("testdata", "testfile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	downloadedData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(testData, downloadedData) {
		t.Fatal("wrong downloaded data")
	}
}

func TestDownloadToFileFailChecksum(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{
		Options: download.Options{
			Checksum:     "d577273f",
			ChecksumHash: crypto.MD5,
		},
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "checksum validation failed") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "checksum validation failed", err)
	}
}

func TestDownloadToFile404(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/invalidfile", tmpFile.Name(), download.FileOptions{
		Options: download.Options{
			Checksum:     "d577273f",
			ChecksumHash: crypto.MD5,
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "received invalid status code") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "received invalid status code", err)
	}
}

func TestDownloadToFileInvalidChecksumHash(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{
		Options: download.Options{
			Checksum:     "d577273f",
			ChecksumHash: crypto.SHA224,
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid hash function") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "invalid hash function", err)
	}
}

type checksum struct {
	checksumFile string
	hash         crypto.Hash
}

var checksumTests = []checksum{
	{"testfile.md5", crypto.MD5},
	{"CHECKSUMS.md5", crypto.MD5},
	{"testfile.sha1", crypto.SHA1},
	{"CHECKSUMS.sha1", crypto.SHA1},
	{"testfile.sha256", crypto.SHA256},
	{"CHECKSUMS.sha256", crypto.SHA256},
	{"testfile.sha512", crypto.SHA512},
	{"CHECKSUMS.sha512", crypto.SHA512},
}

func TestDownloadToFileWithChecksumValidation(t *testing.T) {
	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	for _, chk := range checksumTests {
		func() {
			tmpFile, err := ioutil.TempFile(targetDir, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			defer func() { _ = os.Remove(tmpFile.Name()) }()

			err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{
				Options: download.Options{
					Checksum:     srv.URL + "/" + chk.checksumFile,
					ChecksumHash: chk.hash,
				},
			})
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			testData, err := ioutil.ReadFile(filepath.Join("testdata", "testfile"))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			downloadedData, err := ioutil.ReadFile(tmpFile.Name())
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !bytes.Equal(testData, downloadedData) {
				t.Error("wrong downloaded data")
			}
		}()
	}
}

func TestInvalidURL(t *testing.T) {
	err := download.ToFile("://invalid", "", download.FileOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid src URL") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "invalid src URL", err)
	}
}

func TestNonExistentDestDir(t *testing.T) {
	err := download.ToFile("http://doesnotmatter", filepath.Join("testdata", "nonexistentdir", "somewhere"), download.FileOptions{Mkdirs: download.MkdirNone})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "directory "+filepath.Join("testdata", "nonexistentdir")+" does not exist") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "directory "+filepath.Join("testdata", "nonexistentdir")+" does not exist", err)
	}
}

func TestDownloadToFileSuccessWithRetry(t *testing.T) {
	hfs := http.FileServer(http.Dir("testdata"))
	var srv *httptest.Server
	i := 0
	hf := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if i < 2 {
			i++
			srv.CloseClientConnections()
			return
		}
		hfs.ServeHTTP(w, req)
	})
	srv = httptest.NewServer(hf)
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testData, err := ioutil.ReadFile(filepath.Join("testdata", "testfile"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	downloadedData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(testData, downloadedData) {
		t.Fatal("wrong downloaded data")
	}
}

func TestDownloadToFileFailure(t *testing.T) {
	var srv *httptest.Server
	hf := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		srv.CloseClientConnections()
	})
	srv = httptest.NewServer(hf)
	defer srv.Close()

	targetDir := filepath.Join("testdata", "output")
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.RemoveAll(targetDir) }() // #nosec

	tmpFile, err := ioutil.TempFile(targetDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	err = download.ToFile(srv.URL+"/testfile", tmpFile.Name(), download.FileOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
}
