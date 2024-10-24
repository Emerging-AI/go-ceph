package nfs

import (
	"fmt"

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

type ClusterBeckend struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
}

// ClusterInfo describes an NFS cluster.
type ClusterInfo struct {
	VirtualIP *string          `json:"virtual_ip"`
	Backend   []ClusterBeckend `json:"backend"`
}

func parseClusterInfo(res commands.Response) (map[string]ClusterInfo, error) {
	l := make(map[string]ClusterInfo, 0)
	if err := res.NoStatus().Unmarshal(&l).End(); err != nil {
		return nil, err
	}
	return l, nil
}

// FetchAllClusterInfo returns all information about the Clusters.
//
// Similar To:
//
//	ceph nfs cluster info
func (nfsa *Admin) FetchAllClusterInfo() (map[string]ClusterInfo, error) {

	return nfsa.fetchAllClusterInfo("")
}

// FetchClusterInfo returns information about the specified Cluster.
//
// Similar To:
//
//	ceph nfs cluster info [<cluster_id>]
func (nfsa *Admin) FetchClusterInfo(clusterID string) (*ClusterInfo, error) {
	res, err := nfsa.fetchAllClusterInfo(clusterID)
	if err != nil {
		return nil, err
	}

	cluster, ok := res[clusterID]
	if !ok {
		return nil, fmt.Errorf("cluster %s Not Existed", clusterID)
	}

	return &cluster, err
}

func (nfsa *Admin) fetchAllClusterInfo(clusterID string) (map[string]ClusterInfo, error) {
	m := map[string]string{
		"prefix": "nfs cluster info",
		"format": "json",
	}
	if clusterID != "" {
		m["cluster_id"] = clusterID
	}
	return parseClusterInfo(commands.MarshalMgrCommand(nfsa.conn, m))
}
