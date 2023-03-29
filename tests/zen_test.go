package tests

import (
	"encoding/json"
	"testing"

	"xmind-nodes"
)

func TestXmindZen(t *testing.T) {
	xmindFile, err := xmind_nodes.Load("xmind_zen.xmind")
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := json.MarshalIndent(xmindFile.ExtractAttached(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bytes))
}
