package digitalocean

import "strings"
import "testing"

func TestCloudInit(t *testing.T) {
	if ci := cloudInit("img:1"); !strings.Contains(ci, "img:1") || !strings.Contains(ci, "docker") {
		t.Fatal("bad cloud-init")
	}
}
