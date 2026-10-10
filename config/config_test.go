package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeLeavesCredentialsAlone(t *testing.T) {
	root := t.TempDir()
	credentials := filepath.Join(root, GDDirSuffix, "credentials.json")

	pathGD, firstInit, c, err := Initialize(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !firstInit {
		t.Error("expected the first initialization to be reported as such")
	}
	if want := filepath.Join(root, GDDirSuffix); pathGD != want {
		t.Errorf("gd path: got %q want %q", pathGD, want)
	}
	if _, err := os.Stat(credentials); !os.IsNotExist(err) {
		t.Fatalf("credentials were written before being obtained, stat error: %v", err)
	}

	// Once obtained, they have to survive another initialization
	// for as long as that one doesn't obtain new ones.
	c.ClientId, c.ClientSecret, c.RefreshToken = "client-id", "client-secret", "refresh-token"
	if err := c.Write(); err != nil {
		t.Fatalf("writing the credentials: %v", err)
	}
	want, err := ioutil.ReadFile(credentials)
	if err != nil {
		t.Fatal(err)
	}

	_, firstInit, c, err = Initialize(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if firstInit {
		t.Error("expected the second initialization not to be reported as the first")
	}
	if c.RefreshToken != "" {
		t.Errorf("got a context with the refresh token %q, want a blank one", c.RefreshToken)
	}
	got, err := ioutil.ReadFile(credentials)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("the credentials were overwritten:\ngot  %s\nwant %s", got, want)
	}
}
