// Code generated for package lookup by go-bindata DO NOT EDIT. (@generated)
// sources:
// pkg/health/lookup/static_data/info.yaml
package lookup

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _pkgHealthLookupStatic_dataInfoYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\xd3\x4d\x4e\xc3\x30\x10\x86\xe1\x7d\x4f\xd1\x0b\xb8\xe6\x6f\xd5\x1d\x82\x2e\xd9\x70\x03\xdb\x19\x6a\xc7\xf6\x4c\xe4\x99\xa8\xe9\xed\x51\x4b\x28\x09\x10\x04\x4a\xb8\xc0\xfb\x48\xa3\x6f\x2a\x03\x99\x90\x41\xb6\xab\xf5\xda\x8b\x34\xbc\xd5\x7a\x1f\xc4\xb7\x76\xe3\x28\xeb\x07\xca\xce\xb0\xe8\xd8\x5a\x28\x1e\x4c\x12\x7f\xd4\x36\x91\xd5\x37\xb7\x57\xee\xee\xe5\x5a\xbb\x5c\xe9\x4b\x45\x39\x0f\x2e\xea\xe7\xdd\xfd\xe3\xd3\x6e\x93\xab\x55\x05\x4d\xa2\x63\x06\x9c\x09\x5c\x32\x5f\x05\x64\xc5\x62\xa4\x65\x15\x50\xa0\xa0\x49\xf3\x28\x64\x55\x80\x29\xb5\x12\x08\x7f\xe2\xa0\xfb\x7f\xae\xee\x94\x25\x51\x42\x11\x70\xc2\xa9\x01\x63\x40\x56\x9d\x6a\x52\xbb\x0f\xc8\xba\xee\x54\xf4\x7d\xea\x0c\x66\xc3\x02\xe5\xec\x0d\x83\x63\xc7\x51\x81\xa5\x88\x53\x6b\x5c\x0f\xc8\x62\xd2\xd4\xad\xfe\x0c\xf4\xb9\xb1\xc1\xe0\x0a\x08\x2f\x65\xf4\xb9\xb1\x71\x00\xeb\x89\xe2\x62\xc8\x7b\x6f\xa0\x20\xc8\x81\x4a\x54\x8e\x10\xc1\x7d\xcc\x62\xd6\xce\xa6\xa2\x03\xb7\xa1\xea\x34\x45\x31\x65\xf2\x86\xbf\xb3\x86\xa1\x6f\x95\xb7\x07\x9a\x6d\xf4\x7f\xf8\x49\x78\x0d\x00\x00\xff\xff\x8b\x63\x1c\x00\xd3\x04\x00\x00")

func pkgHealthLookupStatic_dataInfoYamlBytes() ([]byte, error) {
	return bindataRead(
		_pkgHealthLookupStatic_dataInfoYaml,
		"pkg/health/lookup/static_data/info.yaml",
	)
}

func pkgHealthLookupStatic_dataInfoYaml() (*asset, error) {
	bytes, err := pkgHealthLookupStatic_dataInfoYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pkg/health/lookup/static_data/info.yaml", size: 1235, mode: os.FileMode(420), modTime: time.Unix(1602591330, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
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
	"pkg/health/lookup/static_data/info.yaml": pkgHealthLookupStatic_dataInfoYaml,
}

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
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
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
	"pkg": &bintree{nil, map[string]*bintree{
		"health": &bintree{nil, map[string]*bintree{
			"lookup": &bintree{nil, map[string]*bintree{
				"static_data": &bintree{nil, map[string]*bintree{
					"info.yaml": &bintree{pkgHealthLookupStatic_dataInfoYaml, map[string]*bintree{}},
				}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
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
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
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
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
