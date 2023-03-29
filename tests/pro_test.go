package tests

import (
	"encoding/json"
	"testing"

	xmind_nodes "xmind-nodes"
)

func TestXmindPro(t *testing.T) {
	xmindFile, err := xmind_nodes.Load("xmind_pro.xmind")
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := json.MarshalIndent(xmindFile.ExtractAttached(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bytes))
}
