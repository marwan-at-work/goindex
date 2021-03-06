package goindex

import (
	"context"
	"testing"
	"time"
)

func TestGoIndex(t *testing.T) {
	var c Client
	mods, err := c.Get(context.Background(), time.Time{}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(mods) != 10 {
		t.Fatalf("expected 10 mods but got %d", len(mods))
	}
	next, err := mods.Next(context.Background(), &c, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(mods) == 0 {
		t.Fatal("unexpected zero length next")
	}
	last := mods[len(mods)-1]
	first := next[0]
	if last.Path == first.Path && last.Version == first.Version {
		t.Fatalf("expected mods.Next(%+v) to not return a duplicate module", next[0])
	}
}
