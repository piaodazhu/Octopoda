package sdk

import "testing"

func TestSDK(t *testing.T) {
	if err := Init("../octl_test.yaml"); err != nil {
		t.Fatal(err)
	}
	t.Log(NodesInfo(nil))
	res, err := Run("{ls}", []string{"pi4", "pi5"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
