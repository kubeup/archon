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

// +build linux

package download_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	download "github.com/jimmidyson/go-download"
)

func TestNonWritableDestDirCreateSubdir(t *testing.T) {
	_ = os.Chmod(filepath.Join("testdata", "readonlydir"), 0500)
	err := download.ToFile("http://doesnotmatter", filepath.Join("testdata", "readonlydir", "subdir", "somewhere"), download.FileOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to create destination directory") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "failed to create destination directory", err)
	}
}

func TestNonWritableDestDir(t *testing.T) {
	_ = os.Chmod(filepath.Join("testdata", "readonlydir"), 0500)
	err := download.ToFile("http://doesnotmatter", filepath.Join("testdata", "readonlydir", "somewhere"), download.FileOptions{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to create temp file") {
		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "failed to create temp file", err)
	}
}

// func TestNonWritableDestFile(t *testing.T) {
// 	srv := httptest.NewServer(http.FileServer(http.Dir("testdata")))
// 	defer srv.Close()

// 	err := download.ToFile(srv.URL+"/testfile", "testdata/writabledir/readonlyfile", download.FileOptions{})
// 	if err == nil {
// 		t.Fatal("expected error")
// 	}
// 	if !strings.Contains(err.Error(), "failed to rename temp file to destination") {
// 		t.Fatalf("unexpected error, expected to contain: '%s', actual: '%v'", "failed to rename temp file to destination", err)
// 	}
// }
