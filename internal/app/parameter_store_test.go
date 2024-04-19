package app

import (
	"context"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParameterStoreStore(t *testing.T) {
	rules := []ParameterRule{
		{
			Path:  "/foo/v1",
			Level: ParameterLevelStrict,
		},
		{
			Path:  "/bar/",
			Level: ParameterLevelUnder,
		},
		{
			Path:  "/buzz/",
			Level: ParameterLevelAll,
		},
	}

	want := []Parameter{
		{
			Path:  "/foo/v1",
			Value: "this is /foo/v1",
		},
		{
			Path:  "/bar/v1",
			Value: "this is /bar/v1",
		},
		{
			Path:  "/buzz/v1",
			Value: "this is /buzz/v1",
		},
		{
			Path:  "/buzz/a/v2",
			Value: "this is /buzz/a/v2",
		},
		{
			Path:  "/buzz/a/b/v3",
			Value: "this is /buzz/a/b/v3",
		},
	}

	mock := MockSSMConnector{
		data: map[string]string{
			"/foo/v1":      "this is /foo/v1",
			"/foo/v2":      "this is /foo/v2",
			"/bar/v1":      "this is /bar/v1",
			"/bar/a/v2":    "this is /bar/v2",
			"/buzz/v1":     "this is /buzz/v1",
			"/buzz/a/v2":   "this is /buzz/a/v2",
			"/buzz/a/b/v3": "this is /buzz/a/b/v3",
		},
	}

	ctx := context.Background()
	store := NewParameterStore(nil, mock)
	store.Store(ctx, rules)

	sort.Slice(want, func(i, j int) bool {
		return want[i].Path < want[j].Path
	})

	sort.Slice(store.Parameters, func(i, j int) bool {
		return store.Parameters[i].Path < store.Parameters[j].Path
	})

	if diff := cmp.Diff(want, store.Parameters); diff != "" {
		t.Errorf("Store() has diff:\n%s", diff)
	}
}
