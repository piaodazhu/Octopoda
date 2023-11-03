package sdk

import (
	"os"
	"testing"
)

func TestSDK(t *testing.T) {
	if err := Init("../octl_test.yaml"); err != nil {
		t.Fatal(err)
	}
	
	nodesInfo, err := NodesInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nodesInfo)

	nodesStatus, err := NodesStatus(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nodesStatus)

	groups, err := GroupGetAll()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(groups)

	group, err := GroupGet("grp")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(group)

	results, err := Run("{ls}", []string{"pi4", "pi5"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)

	results, err = DistribFile("testfile.txt", "foobar", []string{"pi4", "pi5"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)

	result, err := PullFile("store", "pi4", "foobar/testfile.txt", "localfoobar")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
	
	if _, err = os.Stat("localfoobar/testfile.txt"); os.IsNotExist(err) {
		t.Fatal("PullFile")
	}

	if err := os.RemoveAll("localfoobar"); err != nil {
		t.Fatal(err)
	}
}
