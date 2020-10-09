package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _proto_micro_mall_shop_proto_shop_business_shop_business_swagger_json = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x59\x4d\x6f\xe4\x44\x13\xbe\xcf\xaf\x68\xf5\xfb\x1e\xd1\x66\x09\x88\x43\x6e\xc3\x64\x20\x11\xe4\x43\x33\x13\x09\x84\x22\xab\x63\xd7\x38\xbd\xd8\xdd\x4e\x7f\x84\x0d\x28\x07\x4e\x48\x48\xb0\x8b\x80\x45\x1c\xf6\xc6\x81\x13\x8b\x04\xd2\x4a\x11\x12\x7f\x66\x77\xc8\xfe\x0b\xd4\xfe\x1a\xbb\xc7\x9e\xcc\xd8\x09\x49\xb4\xb1\x14\x29\xee\xee\x7a\xba\xaa\xba\x9e\x72\x55\xcf\x17\x1d\x84\xb0\xfc\x8c\xf8\x3e\x08\xbc\x86\xf0\xea\xbd\xfb\xf8\x0d\x33\x46\xd9\x98\xe3\x35\x64\xe6\x11\xc2\x8a\xaa\x00\xcc\x7c\x24\xb8\xe2\x2b\x21\x75\x05\x77\x42\x12\x04\x8e\x3c\xe4\x91\x93\x8c\xc6\xff\x1e\x68\x49\x19\x48\x59\x7e\xbb\x17\xaf\x88\x91\x11\xc2\xc7\x20\x24\xe5\xcc\xe0\xa5\xff\x22\xc6\x15\x92\xa0\x70\x07\xa1\xd3\x78\x7f\x97\x33\xa9\x43\x90\x78\x0d\x7d\x92\x48\x91\x28\x0a\xa8\x4b\x14\xe5\x6c\xe5\x81\xe4\xcc\xac\xdd\x8f\xd7\x46\x82\x7b\xda\x5d\x70\x2d\x51\x87\x72\x6a\xd8\xca\xf1\x9b\x15\x7a\xaf\x18\xf9\x93\x7c\x95\x11\xe3\x52\x15\xde\x8d\xd3\x74\x18\x12\x61\x16\xe1\xc9\xa3\xc7\x2f\xce\x7e\x79\x79\xf6\xd3\xab\xef\xcf\xfe\xf9\xe1\x8f\xf3\x67\xcf\x27\x4f\x1f\x4f\x9e\xfc\x9c\xda\x1b\x2f\xe7\x11\x88\x58\x9f\x4d\xcf\x88\x0c\x0f\x79\xf4\x6e\xba\xe5\x10\xc4\x31\x75\xc1\x31\x63\xdd\x78\xe3\x82\x9c\x00\x19\x71\x26\x41\x96\x76\x47\x08\xaf\xde\xbf\x6f\x0d\x21\x84\x3d\x90\xae\xa0\x91\x4a\xbd\xdb\x45\x52\xbb\x2e\x48\x39\xd6\x01\xca\x90\xee\x15\xe0\x13\x4b\xdc\x43\x08\xc9\x0c\x18\x42\xf8\xff\x02\xc6\x06\xe7\x7f\x2b\x1e\x8c\x29\xa3\x06\xd7\x3a\xd9\x5c\xe9\x41\x0a\x8f\x4b\x20\xa7\x85\xb7\xd3\xe2\xbe\xd8\x83\x31\xd1\x81\xba\xd8\x06\x86\x34\x83\x87\x11\xb8\x0a\x3c\x04\x42\x70\x91\x9b\xd2\xd6\x12\xa1\x99\xa2\x21\xf4\x0d\xe8\x1c\xbd\x3b\x15\x16\xe0\x88\x08\x12\x82\x02\x31\x0d\xbb\xe4\xb1\xcc\x61\x24\x8c\x99\x73\xc0\xbd\x13\x5b\x5f\xca\xea\x66\x04\x1c\x69\x2a\xc0\xc4\x8a\x12\x1a\xae\xe8\xc4\x8e\x34\x48\xb5\x88\xe1\xfb\x05\xc3\x15\xf1\x6d\x93\xab\xc2\x79\x0a\xbb\xdf\x29\xc2\xa5\x3e\xac\x63\x5e\x29\xef\x20\x84\x7d\xa8\xe7\xdd\xf9\xb7\xcf\x5f\x3e\x7a\x92\xf0\x6e\xf2\xe3\xef\x93\x6f\x7e\x5b\x92\x71\xef\x83\x32\xc3\x9b\x66\xcf\x5b\xc4\xb9\x82\xda\x77\xac\x4b\x9f\x1a\xd6\xc5\x7e\xa3\x9e\xac\x66\xde\x91\x06\x31\x8f\x7a\x63\x12\x48\x9b\x7b\xea\x24\x8a\x91\x89\x10\x64\x96\xd0\x0a\x42\x3b\x6a\x4a\x52\x52\x09\xca\x7c\x4b\x0c\x21\x3c\xe6\x22\x24\xe6\x58\x30\x65\xea\x9d\xb7\x2d\xb7\x58\xbb\xb8\x3c\x08\xc0\x35\xce\x7c\x2f\x17\x0b\x75\xa0\x28\xbe\x66\xf2\x86\x44\x81\xa0\x24\x68\x44\xe0\x26\x9f\xcc\x94\x09\x5b\xd9\xbe\xb7\x8f\xc4\x99\xea\x77\x44\x4e\x9f\xf9\x44\xbe\x7c\x1e\x57\x32\x72\x1e\x1f\xaf\x83\x58\x51\x00\x9e\x0f\xa5\x82\x54\xd7\xd3\x2a\x21\xd4\xf9\x9f\xbf\x4e\xbe\xfe\xeb\xc5\xdf\x4f\xcf\x9f\x7d\xf9\xea\xab\xef\x1a\x14\xa3\xbb\xc9\xae\xb7\x88\x54\x53\xad\xef\xf8\x94\x3e\xb7\xa1\x1c\xcd\x8e\xec\x5a\xea\xd1\xbc\xe5\x2c\xa8\x38\x6d\x10\xe3\xc6\xf5\x40\x8f\xbb\xac\xd4\x0f\x66\xd9\x83\x1f\x3c\x00\x57\xe5\xee\x32\xcb\x23\x10\x8a\x5a\x0c\x89\xd7\x3b\x5a\x04\x36\x6f\xac\x2c\x54\x79\xb4\xc7\x24\xd0\x70\x81\x60\x29\x96\xa7\xc9\xeb\xe0\x44\x15\x0c\x3f\xad\x4c\x39\xa5\xa8\x6b\x61\x22\x58\x00\x0b\xdb\xe7\x72\xaf\xd6\x3c\xca\x14\xf8\x20\xea\xec\xa3\x4c\xbd\xb5\x5a\x8d\x1a\x82\x94\xc4\xbf\xc8\x6f\x95\xa2\x1e\x28\x42\x83\x99\x1c\x57\x5f\xf9\xd5\xd4\x7d\x35\xe1\x5f\x0c\xa9\x6a\x5e\x57\x9e\x53\x89\x35\x3d\x4e\xd9\x28\xd1\x67\xe6\xc0\xac\x90\xc0\xc0\x74\x58\x22\x08\xee\x6d\x7f\x5c\xcc\xea\x7b\xc3\xf5\xe2\xeb\xc6\x07\xa5\xd7\xfe\xde\x20\x53\x33\xa7\x5d\x21\x5d\xc6\x60\xf3\x14\x0d\x43\xce\xf2\x6c\xdc\x22\xbe\xaa\xc2\x64\x81\xfc\x32\x00\xd5\x33\x92\xd5\x41\x22\xfd\x85\x03\xe4\xe2\x43\xa9\xea\xcb\x5a\x19\x6c\x5c\xd7\xc0\x64\xcb\xe7\x95\x96\x9b\x3e\xdb\x09\xa8\xb4\xbf\x78\x97\x16\xe5\x33\x49\x3e\xee\xb2\x9b\xc6\x7b\x5d\xb5\x7c\x53\xdd\x5b\xd1\x09\x2d\xe3\xab\xbc\xa1\x59\xc6\x47\x3b\x59\x55\xd7\x22\x31\x0c\xfa\xdd\x51\xbf\x94\x1b\x76\xd7\xad\x91\xf5\xfe\x87\xfd\xf2\x48\x77\x6f\x7d\x73\x34\x3f\x47\x24\xb8\xf5\xba\x67\x2c\x6d\xa4\xf5\x70\xaf\xd7\xeb\x0f\x87\xa5\xa4\x35\x18\xec\x0c\xca\x39\xae\x3f\x70\xb6\x77\x46\x4e\xff\xa3\xcd\xe1\x68\x66\x66\x66\x74\xab\x3f\xe8\x6d\x74\xb7\x47\xd5\x32\xf9\xec\xcc\xcc\x70\x63\x67\xb7\x5a\x26\x9e\x49\x46\xe7\xb9\x2a\x33\xa6\xde\x57\x33\x17\x78\x2d\x48\x90\x77\x02\x8e\x2a\x07\xcd\xa2\xd1\x5a\x0e\xba\x4a\x2e\xf0\xc8\xd1\xd4\x6b\x56\xc0\x58\xdd\x97\x05\x4b\xa3\x26\x9f\xf7\xac\x95\xbc\x6c\x85\x42\x10\xee\x21\x61\xea\x2a\xb0\x19\x75\x3f\x75\xd2\xda\x7d\x69\x83\xc7\x3a\x08\x1a\x0b\x0b\xf0\xa9\x54\x20\x1c\xe2\x79\x8d\xca\xbb\x2c\x52\xda\x03\x04\xd4\x85\x72\xd2\x5f\x18\x43\x91\x87\x8e\x4b\x84\xe7\x30\xde\x4a\x05\xd3\x0d\x36\x8a\x39\xee\x52\x12\x38\xae\x00\x8f\x2a\x67\x5e\xc5\x3b\x0f\x85\x0b\x9f\x30\xfa\x79\x42\xd8\xa5\x40\x2e\xfe\x7a\xcc\xfe\x90\x73\x53\xbf\xad\x97\xca\xdf\xc5\x1c\xb3\x69\xfd\x2a\xb1\xac\x3f\x6e\x63\xca\x69\x95\x35\x62\x83\xaf\x20\x42\xb7\x2a\x6e\x98\x5f\x87\xc3\xb8\xcb\xff\x77\xf9\xbf\x4d\xfe\xff\x0f\xa8\x59\xbe\xc4\x6b\x55\x93\xbe\x1e\xf5\x22\x09\xb9\x66\xb5\xfd\xf7\xfc\x1b\x33\xda\xb8\x60\xcf\x6f\x8e\x9a\x9e\xef\x4d\xaf\x0e\x14\x51\x97\x75\x9d\x58\x76\x4c\xc7\xfc\x9d\x76\xfe\x0d\x00\x00\xff\xff\x1a\x8e\x90\xdf\x33\x25\x00\x00")

func proto_micro_mall_shop_proto_shop_business_shop_business_swagger_json() ([]byte, error) {
	return bindata_read(
		_proto_micro_mall_shop_proto_shop_business_shop_business_swagger_json,
		"proto/micro_mall_shop_proto/shop_business/shop_business.swagger.json",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"proto/micro_mall_shop_proto/shop_business/shop_business.swagger.json": proto_micro_mall_shop_proto_shop_business_shop_business_swagger_json,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func     func() ([]byte, error)
	Children map[string]*_bintree_t
}

var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"proto": &_bintree_t{nil, map[string]*_bintree_t{
		"micro_mall_shop_proto": &_bintree_t{nil, map[string]*_bintree_t{
			"shop_business": &_bintree_t{nil, map[string]*_bintree_t{
				"shop_business.swagger.json": &_bintree_t{proto_micro_mall_shop_proto_shop_business_shop_business_swagger_json, map[string]*_bintree_t{}},
			}},
		}},
	}},
}}