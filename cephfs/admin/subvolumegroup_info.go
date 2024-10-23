package admin

// SubVolumeGroupInfo reports various informational values about a subvolumeGroup.
type SubVolumeGroupInfo struct {
	Uid          int       `json:"uid"`
	Gid          int       `json:"gid"`
	Mode         int       `json:"mode"`
	MonAddrs     []string  `json:"mon_addrs"`
	BytesPercent string    `json:"bytes_pcent"`
	BytesUsed    ByteCount `json:"bytes_used"`
	BytesQuota   QuotaSize `json:"-"`
	DataPool     string    `json:"data_pool"`
	Atime        TimeStamp `json:"atime"`
	Mtime        TimeStamp `json:"mtime"`
	Ctime        TimeStamp `json:"ctime"`
	CreatedAt    TimeStamp `json:"created_at"`
}

type subVolumeGroupInfoWrapper struct {
	SubVolumeGroupInfo
	VBytesQuota *quotaSizePlaceholder `json:"bytes_quota"`
}

func parseSubVolumeGroupInfo(res response) (*SubVolumeGroupInfo, error) {
	var info subVolumeGroupInfoWrapper
	if err := res.NoStatus().Unmarshal(&info).End(); err != nil {
		return nil, err
	}
	if info.VBytesQuota != nil {
		info.BytesQuota = info.VBytesQuota.Value
	}
	return &info.SubVolumeGroupInfo, nil
}

// SubVolumeInfo returns information about the specified subvolume.
//
// Similar To:
//
//	ceph fs subvolumegroup info <volume> <group_name>
func (fsa *FSAdmin) FetchSubVolumeGroupInfo(volume, name string) (*SubVolumeGroupInfo, error) {
	m := map[string]string{
		"prefix":     "fs subvolumegroup info",
		"vol_name":   volume,
		"group_name": name,
		"format":     "json",
	}

	return parseSubVolumeGroupInfo(fsa.marshalMgrCommand(m))
}
