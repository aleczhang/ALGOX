package settings

import "testing"

func TestSaveLoadAndRefuse(t *testing.T) {
	t.Setenv("ALGOX_HOME", t.TempDir())
	if s, err := Load(); err != nil || s != Default() {
		t.Fatalf("fresh: %+v %v", s, err)
	}
	want := Settings{Theme: "dark", Lang: "zh"}
	if err := Save(want); err != nil {
		t.Fatal(err)
	}
	if s, _ := Load(); s != want {
		t.Fatalf("got %+v", s)
	}
	if err := Save(Settings{Theme: "pink", Lang: "zh"}); err == nil {
		t.Fatal("unknown theme saved")
	}
	if s, _ := Load(); s != want {
		t.Fatalf("a refused save changed the file: %+v", s)
	}
}
