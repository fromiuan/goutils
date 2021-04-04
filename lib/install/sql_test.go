package install

import (
	"strings"
	"testing"
)

func Test_Exec(t *testing.T) {
	conn, err := NewInstall("127.0.0.1", 3306, "root", "admin", "demo")
	if err != nil {
		t.Fatal(err)
	}
	conn.ExecReder(strings.NewReader("SELECT * FROM user"))

}

func Test_NewInstall(t *testing.T) {
	_, err := NewInstall("127.0.0.1", 3306, "root", "admin", "demo")
	if err != nil {
		t.Fatal(err)
	}
}
