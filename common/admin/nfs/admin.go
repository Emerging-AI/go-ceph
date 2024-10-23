//go:build !(nautilus || octopus)
// +build !nautilus,!octopus

package nfs

import (
	"strings"

	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/internal/commands"
)

// Admin is used to administer ceph nfs features.
type Admin struct {
	conn ccom.RadosCommander
}

// NewFromConn creates an new management object from a preexisting
// rados connection. The existing connection can be rados.Conn or any
// type implementing the RadosCommander interface.
func NewFromConn(conn ccom.RadosCommander) *Admin {
	return &Admin{conn}
}

// type listNamedResult struct {
// 	Name string `json:"name"`
// }

// func parseListNames(res commands.Response) ([]string, error) {
// 	var r []listNamedResult
// 	if err := res.NoStatus().Unmarshal(&r).End(); err != nil {
// 		return nil, err
// 	}
// 	vl := make([]string, len(r))
// 	for i := range r {
// 		vl[i] = r[i].Name
// 	}
// 	return vl, nil
// }

// parseClusterNameResponse returns a cleaned up cluster name from requests that get a list of names
// unless an error is encountered, then an error is returned.
func parseClusterNameResponse(res commands.Response) ([]string, error) {
	if res2 := res.NoStatus(); !res2.Ok() {
		return nil, res.End()
	}
	b := res.Body()
	for len(b) >= 1 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}

	clNames := strings.Split(string(b), "\\n")

	cleanedClNames := make([]string, 0)
	for _, item := range clNames {
		if item != "" {
			cleanedClNames = append(cleanedClNames, item)
		}
	}

	return cleanedClNames, nil
}
