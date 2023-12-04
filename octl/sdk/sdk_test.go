package sdk

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSDK(t *testing.T) {
	if err := Init("../octl_test.yaml"); err != nil {
		t.Fatal(err)
	}

	nodesInfo, err := NodeInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nodesInfo)

	nodesStatus, err := NodeStatus(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nodesStatus)

	results, err := RunCommand("uname -a", false, []string{"pi4", "pi5"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)

	results, err = UploadFile("testfile.txt", "@fstore/foobar", true, []string{"pi4", "pi5"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)

	result, err := DownloadFile("@fstore/foobar/testfile.txt", "localfoobar", "pi4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)

	if _, err := os.Stat("localfoobar/testfile.txt"); os.IsNotExist(err) {
		t.Fatal("PullFile")
	}

	if err := os.RemoveAll("localfoobar"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*10)
	defer cancel()
	logList, err := Apply(ctx, "../s1", "stop", "byGoCli")
	t.Log(strings.Join(logList, "\n"))
	t.Log(err)
}
