package vmd

import (
	"os"
	"testing"
)

func TestDecodeMotion(t *testing.T) {
	f, err := os.Open("testdata/motion.vmd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if p.Header.Version < 1 || p.Header.Version > 2 {
		t.Fatal(err)
	}
}

func TestDecodeCamera(t *testing.T) {
	f, err := os.Open("testdata/camera.vmd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	p, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if p.Header.Version < 1 || p.Header.Version > 2 {
		t.Fatal(err)
	}
}
