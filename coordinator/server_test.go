package main

import (
	"context"
	"github.com/google/go-cmp/cmp"
	"github.com/wangwalton/gocrawler/contracts"
	"testing"
)

func TestHostnameCoordinatorServer_AddHostnames(t *testing.T) {
	s := newServer()

	hostname := &contracts.HostnamePaths{
		Hostname:                "google.com",
		NumReqsPerMinuteAllowed: 0,
		Paths: map[string]*contracts.Empty{
			"maps":  {},
			"drive": {},
		},
	}

	req := &contracts.MultipleHostnamePaths{
		Urls: []*contracts.HostnamePaths{hostname},
	}

	_, err := s.AddHostnames(context.Background(), req)
	if err != nil {
		t.Errorf("AddHostname(%v) got unexpected error", err)
	}

	if len(s.hostnameQueue) != 1 {
		t.Errorf("AddHostname wanted s.hostnameQueue to be len 1, got %d", len(s.hostnameQueue))
	}

	if len(s.hostnameMap) != 1 {
		t.Errorf("AddHostname wanted s.hostnameMap to be len 1, got %d", len(s.hostnameMap))
	}

	val, ok := s.hostnameMap[hostname.Hostname]
	if !ok {
		t.Errorf("AddHostname expented %s to be in map ", hostname.Hostname)
	}

	if diff := cmp.Diff(*hostname, *val); diff != "" {
		t.Errorf("AddHostname() mismatch (-want +got):\n%s", diff)
	}
}
