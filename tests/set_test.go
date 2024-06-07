package tests

import (
	"bytes"
	"fmt"
	"testing"
)

// func runTestServer(wg *sync.WaitGroup) {
// 	config.IOBufferLength = 16
// 	config.Port = 8379
// 	wg.Add(1)
// 	server.RunAsyncTCPServer(wg)
// }

func TestSet(t *testing.T) {
	// var wg sync.WaitGroup
	// go runTestServer(&wg)

	var b []byte
	var buf = bytes.NewBuffer(b)
	conn := getLocalConnection()
	for i := 1; i < 2; i++ {
		buf.WriteByte('a')
		cmd := fmt.Sprintf("SET k%d %s", i, buf.String())
		fireCommand(conn, cmd)
	}

	for i := 1; i < 2; i++ {
		cmd := fmt.Sprintf("GET k%d", i)
		v := fireCommand(conn, cmd)
		if len(v.(string)) != i {
			t.Fail()
		}
	}
	fireCommand(conn, "ABORT")
	// wg.Wait()
}
