package nfs

import (
	"github.com/ceph/go-ceph/internal/commands"
)

// ListClusters will return a list of clusters.
//
// Similar To:
//
//	ceph nfs cluster ls
func (nfsa *Admin) ListClusters() ([]string, error) {
	m := map[string]string{
		"prefix": "nfs cluster ls",
		"format": "json",
	}
	res := commands.MarshalMgrCommand(nfsa.conn, m)

	return parseClusterNameResponse(res)
}
