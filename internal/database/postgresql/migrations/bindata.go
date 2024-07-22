// Code generated by go-bindata. DO NOT EDIT.
// sources:
// V1.0__create_schema.sql (607B)

package migrations

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _v10__create_schemaSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x92\xcb\x4a\x03\x31\x14\x86\xf7\xf3\x14\xff\x72\x02\x22\x15\x15\x17\x5d\xc5\x36\x65\x82\x73\x29\x99\x33\xda\x76\x33\x84\x26\x4a\xb0\x9d\x96\xb9\x78\x79\x7b\x69\xa6\x45\x2b\x62\xfb\x2f\x42\xc2\xf9\x3e\xce\x09\xc9\x48\x09\x4e\x02\xf9\x28\x12\x09\x87\x9c\x20\xcd\x08\x62\x26\x73\xca\xd1\x35\xb6\x6e\xc0\x0b\x8a\x32\x25\x17\x9c\x64\x96\x42\x9b\xb5\xab\x86\x41\xb0\xf7\x88\xdf\xc7\xa2\x07\x2f\xfd\x1a\x84\x01\x00\x38\x83\x43\x8a\x42\x8e\x31\x16\x13\x5e\xc4\x84\xae\x73\xa6\x7c\xb1\x95\xad\x75\x6b\xcb\xb7\x9b\x90\x5d\x78\xfe\xd9\xd5\x4d\x5b\x56\x7a\x6d\xf1\xc8\xd5\x28\xe2\x2a\xbc\xbb\x65\x38\xca\x6e\xb2\xb4\x88\xe3\xde\x58\xe9\x83\x70\xae\x51\xb9\xe5\x6b\x2f\x9c\x6b\x6c\x75\xd3\xbc\x6f\xea\xdd\x5d\x48\xcc\x08\x7f\xe5\xd8\xb0\x6b\xed\x56\x7d\xe1\xd0\xe3\x6a\x30\x60\xff\x18\xcb\x4d\x57\xb5\xf5\xe7\x6e\xeb\xf1\x6b\x76\xaa\x47\xb7\x35\xba\xb5\xa6\xd4\x2d\x48\x26\x22\x27\x9e\x4c\xf1\x24\x29\xf2\x47\x2c\xb2\x54\xfc\x32\x8c\x5d\xd9\x93\x86\xa7\x3d\x3e\x55\x32\xe1\x6a\x8e\x07\x31\x47\xe8\x0c\x0b\xd8\xf7\x83\xcb\x74\x2c\x66\x70\xe6\xa3\xdc\x0f\xee\x8d\x2c\xfd\xf9\x07\x10\xee\x6b\x6c\xf8\x15\x00\x00\xff\xff\xb9\x32\x56\x19\x5f\x02\x00\x00")

func v10__create_schemaSqlBytes() ([]byte, error) {
	return bindataRead(
		_v10__create_schemaSql,
		"V1.0__create_schema.sql",
	)
}

func v10__create_schemaSql() (*asset, error) {
	bytes, err := v10__create_schemaSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "V1.0__create_schema.sql", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd9, 0x26, 0xc2, 0x79, 0xfe, 0x9b, 0x9a, 0x5b, 0x98, 0x15, 0x3b, 0xf2, 0xd5, 0x16, 0xa1, 0xcd, 0x4f, 0x5a, 0x54, 0xb2, 0x63, 0xb, 0xd1, 0x81, 0xd7, 0x99, 0x5f, 0x90, 0x29, 0xa8, 0xbf, 0x81}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"V1.0__create_schema.sql": v10__create_schemaSql,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"V1.0__create_schema.sql": {v10__create_schemaSql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory.
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = os.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}
