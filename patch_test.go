package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

func JsonEqual(a, b []byte) bool {
	var va, vb interface{}
	if err := json.Unmarshal(a, &va); err != nil {
		fmt.Printf("json.Unmarshal(%s) error\n", a)
		return false
	}
	if err := json.Unmarshal(b, &vb); err != nil {
		fmt.Printf("json.Unmarshal(%s) error\n", b)
		return false
	}

	var ma, mb map[string]interface{}
	ma, ok := va.(map[string]interface{})
	if !ok {
		fmt.Printf("%v not map\n", va)
		return false
	}
	mb, ok = vb.(map[string]interface{})
	if !ok {
		fmt.Printf("%v not map\n", vb)
		return false
	}
	/*
		fmt.Println("\na map:")
		for k,v := range ma {
			fmt.Printf("[%s]: %v\n", k, v)
		}
		fmt.Println("\nb map:")
		for k,v := range mb {
			fmt.Printf("[%s]: %v\n", k, v)
		}
	*/
	return reflect.DeepEqual(ma, mb)
}

func TomlEqual(a, b []byte) bool {
	var va, vb interface{}
	if err := toml.Unmarshal(a, &va); err != nil {
		fmt.Printf("toml.Unmarshal(%s) error\n", a)
		return false
	}
	if err := toml.Unmarshal(b, &vb); err != nil {
		fmt.Printf("toml.Unmarshal(%s) error\n", b)
		return false
	}

	var ma, mb map[string]interface{}
	ma, ok := va.(map[string]interface{})
	if !ok {
		fmt.Printf("%v not map\n", va)
		return false
	}
	mb, ok = vb.(map[string]interface{})
	if !ok {
		fmt.Printf("%v not map\n", vb)
		return false
	}
	/*
		fmt.Println("\na map:")
		for k,v := range ma {
			fmt.Printf("[%s]: %v\n", k, v)
		}
		fmt.Println("\nb map:")
		for k,v := range mb {
			fmt.Printf("[%s]: %v\n", k, v)
		}
	*/
	return reflect.DeepEqual(ma, mb)
}

func TestTomlEqual(t *testing.T) {
	var a = `
root = "/var/lib/containerd"
state = "/run/containerd"
`

	var a1 = `
root = "/data/containerd"
`
	var a2 = `
state = "/run/containerd"
root = "/var/lib/containerd"
`

	var b = `
root = "/var/lib/containerd"
[plugins]
  [plugins.cri]
    stream_server_port = "8888"
    [plugins.cri.cni]
      bin_dir = "/usr/bin"

`
	var b1 = `
root = "/var/lib/containerd"
[plugins]
  [plugins.cri.cni]
      bin_dir = "/usr/bin"
  [plugins.cri]
    stream_server_port = "8888"
`

	var b2 = `
root = "/var/lib/containerd"
[plugins.cri]
    stream_server_port = "8888"
[plugins.cri.cni]
    bin_dir = "/usr/bin"

`

	assert.True(t, TomlEqual([]byte(a), []byte(a)))
	assert.False(t, TomlEqual([]byte(a), []byte(a1)))
	assert.True(t, TomlEqual([]byte(a), []byte(a2)))
	assert.True(t, TomlEqual([]byte(b), []byte(b1)))
	assert.True(t, TomlEqual([]byte(b1), []byte(b2)))
}

func TestPatchToml(t *testing.T) {
	result, err := PatchConfigFile("test/containerd-config.toml",
		"test/containerd-patch.toml", "auto")
	assert.Nil(t, err)

	expectedToml, err := ioutil.ReadFile("test/containerd-merge.toml")
	assert.Nil(t, err)

	assert.True(t, TomlEqual(result, expectedToml))

}

func TestPatchJson(t *testing.T) {
	result, err := PatchConfigFile("test/docker-daemon.json",
		"test/docker-patch.json", "auto")
	assert.Nil(t, err)

	expectedJson, err := ioutil.ReadFile("test/docker-merge.json")
	assert.Nil(t, err)

	assert.True(t, JsonEqual(result, expectedJson))
}
