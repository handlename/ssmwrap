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

func TestParameterStoreRetrieve(t *testing.T) {
	paramAttrs := map[string]string{
		"/foo/v1":   "this is /foo/v1",
		"/bar/v1":   "this is /bar/v1",
		"/bar/v2":   "this is /bar/v2",
		"/bar/a/v3": "this is /bar/a/v3",
	}

	params := []Parameter{}
	for p, v := range paramAttrs {
		params = append(params, Parameter{
			Path:  p,
			Value: v,
		})
	}

	store := ParameterStore{
		Parameters: params,
	}

	tests := []struct {
		title string
		path  string
		level ParameterLevel
		want  []Parameter
	}{
		{
			title: "strict",
			path:  "/bar/a/v3",
			level: ParameterLevelStrict,
			want: []Parameter{
				{Path: "/bar/a/v3", Value: paramAttrs["/bar/a/v3"]},
			},
		},
		{
			title: "under",
			path:  "/bar/",
			level: ParameterLevelUnder,
			want: []Parameter{
				{Path: "/bar/v1", Value: paramAttrs["/bar/v1"]},
				{Path: "/bar/v2", Value: paramAttrs["/bar/v2"]},
			},
		},
		{
			title: "all",
			path:  "/bar/",
			level: ParameterLevelAll,
			want: []Parameter{
				{Path: "/bar/v1", Value: paramAttrs["/bar/v1"]},
				{Path: "/bar/v2", Value: paramAttrs["/bar/v2"]},
				{Path: "/bar/a/v3", Value: paramAttrs["/bar/a/v3"]},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got, _ := store.Retrieve(tt.path, tt.level)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Retrieve() has diff:\n%s", diff)
			}
		})
	}
}
