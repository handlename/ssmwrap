package app

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/google/go-cmp/cmp"
	"github.com/samber/lo"
)

type MockSSMConnector struct {
	data map[string]string
}

func (c MockSSMConnector) fetchParametersByPaths(ctx context.Context, client *ssm.Client, paths []string, recursive bool) (map[string]string, error) {
	ret := map[string]string{}
	dataKeys := lo.Keys(c.data)

	for _, path := range paths {
		keys := lo.Filter(dataKeys, func(key string, _ int) bool {
			if !strings.HasPrefix(key, path) {
				return false
			}

			if !recursive {
				rest := strings.Replace(key, path, "", 1)
				if strings.Contains(rest, "/") {
					return false
				}
			}

			return true
		})

		if len(keys) == 0 {
			continue
		}

		for _, key := range keys {
			ret[key] = c.data[key]
		}
	}

	return ret, nil
}

func (c MockSSMConnector) fetchParametersByNames(ctx context.Context, client *ssm.Client, names []string) (map[string]string, error) {
	ret := map[string]string{}

	for _, name := range names {
		v, ok := c.data[name]
		if ok {
			ret[name] = v
		}
	}

	return ret, nil
}

func TestMockSSMConnectorFetchParametersByPaths(t *testing.T) {
	mock := MockSSMConnector{
		data: map[string]string{
			"/foo/bar":          "this is /foo/bar",
			"/bar/v1":           "this is /bar/v1",
			"/bar/v2":           "this is /bar/v2",
			"/buzz/v1":          "this is /buzz/v1",
			"/buzz/qux/v2":      "this is /buzz/qux/v2",
			"/buzz/qux/quux/v3": "this is /buzz/qux/quux/v3",
		},
	}

	test := []struct {
		title     string
		paths     []string
		recursive bool
		want      map[string]string
	}{
		{
			title:     "just under",
			paths:     []string{"/bar/", "/buzz/"},
			recursive: false,
			want: map[string]string{
				"/bar/v1":  "this is /bar/v1",
				"/bar/v2":  "this is /bar/v2",
				"/buzz/v1": "this is /buzz/v1",
			},
		},
		{
			title:     "recursively",
			paths:     []string{"/bar/", "/buzz/"},
			recursive: true,
			want: map[string]string{
				"/bar/v1":           "this is /bar/v1",
				"/bar/v2":           "this is /bar/v2",
				"/buzz/v1":          "this is /buzz/v1",
				"/buzz/qux/v2":      "this is /buzz/qux/v2",
				"/buzz/qux/quux/v3": "this is /buzz/qux/quux/v3",
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.title, func(t *testing.T) {
			got, err := mock.fetchParametersByPaths(context.Background(), nil, tt.paths, tt.recursive)
			if err != nil {
				t.Errorf("fetchParametersByPaths() error = %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("fetchParametersByPaths() has diff:\n%s", diff)
			}
		})
	}
}

func TestMockSSMConnectorFetchParametersByNames(t *testing.T) {
	mock := MockSSMConnector{
		data: map[string]string{
			"/foo/v1": "this is /foo/v1",
			"/bar/v2": "this is /bar/v2",
		},
	}

	test := []struct {
		title     string
		names     []string
		recursive bool
		want      map[string]string
	}{
		{
			title:     "success",
			names:     []string{"/foo/v1", "/bar/v2"},
			recursive: false,
			want: map[string]string{
				"/foo/v1": "this is /foo/v1",
				"/bar/v2": "this is /bar/v2",
			},
		},
		{
			title:     "no result",
			names:     []string{"/unknown/value"},
			recursive: true,
			want:      map[string]string{},
		},
	}

	for _, tt := range test {
		t.Run(tt.title, func(t *testing.T) {
			got, err := mock.fetchParametersByPaths(context.Background(), nil, tt.names, tt.recursive)
			if err != nil {
				t.Errorf("fetchParametersByPaths() error = %v", err)
				return
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("fetchParametersByNames() has diff:\n%s", diff)
			}
		})
	}
}
