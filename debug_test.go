package debug

import "strings"
import "testing"
import "bytes"

func assertContains(t *testing.T, str, substr string) {
	if !strings.Contains(str, substr) {
		t.Fatalf("expected %q to contain %q", str, substr)
	}
}

func assertNotContains(t *testing.T, str, substr string) {
	if strings.Contains(str, substr) {
		t.Fatalf("expected %q to not contain %q", str, substr)
	}
}

func TestDefault(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Writer = buf

	debug := Debug("foo")
	debug("something")
	debug("here")
	debug("whoop")

	if buf.Len() != 0 {
		t.Fatalf("buffer should be empty")
	}
}

func TestEnable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Writer = buf

	Enable("foo")

	debug := Debug("foo")
	debug("something")
	debug("here")
	debug("whoop")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "something")
	assertContains(t, str, "here")
	assertContains(t, str, "whoop")
}

func TestMultipleOneEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Writer = buf

	Enable("foo")

	foo := Debug("foo")
	foo("foo")

	bar := Debug("bar")
	bar("bar")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "foo")
	assertNotContains(t, str, "bar")
}

func TestMultipleEnabled(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Writer = buf

	Enable("foo,bar")

	foo := Debug("foo")
	foo("foo")

	bar := Debug("bar")
	bar("bar")

	if buf.Len() == 0 {
		t.Fatalf("buffer should have output")
	}

	str := string(buf.Bytes())
	assertContains(t, str, "foo")
	assertContains(t, str, "bar")
}

func TestEnableDisable(t *testing.T) {
	var b []byte
	buf := bytes.NewBuffer(b)
	Writer = buf

	Enable("foo,bar")
	Disable()

	foo := Debug("foo")
	foo("foo")

	bar := Debug("bar")
	bar("bar")

	if buf.Len() != 0 {
		t.Fatalf("buffer should not have output")
	}
}

func BenchmarkDisabled(b *testing.B) {
	debug := Debug("something")
	for i := 0; i < b.N; i++ {
		debug("stuff")
	}
}

func BenchmarkNonMatch(b *testing.B) {
	debug := Debug("something")
	Enable("nonmatch")
	for i := 0; i < b.N; i++ {
		debug("stuff")
	}
}
