package admin

// SubVolumeGroupInfo reports various informational values about a subvolumegroup.
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

// FetchSubVolumeGroupInfo returns information about the specified subvolumegroup.
//
// Similar To:
//
//	quincy: ceph fs subvolumegroup info <volume> <group_name>
func (fsa *FSAdmin) FetchSubVolumeGroupInfo(volume, name string) (*SubVolumeGroupInfo, error) {
	m := map[string]string{
		"prefix":     "fs subvolumegroup info",
		"vol_name":   volume,
		"group_name": name,
		"format":     "json",
	}

	return parseSubVolumeGroupInfo(fsa.marshalMgrCommand(m))
}

type subVolumeGroupResizeFields struct {
	Prefix    string `json:"prefix"`
	Format    string `json:"format"`
	VolName   string `json:"vol_name"`
	GroupName string `json:"group_name,omitempty"`
	NewSize   string `json:"new_size"`
	NoShrink  bool   `json:"no_shrink"`
}

// SubVolumeGroupResizeResult reports the size values returned by the
// ResizeSubVolume function, as reported by Ceph.
type SubVolumeGroupResizeResult struct {
	BytesUsed    ByteCount `json:"bytes_used"`
	BytesQuota   ByteCount `json:"bytes_quota"`
	BytesPercent string    `json:"bytes_pcent"`
}

// ResizeSubVolumeGroup will delete a subvolume group in a volume.
// Similar To:
//
//	ceph fs subvolumegroup resize <vol_name> <group_name> <new_size> [-no_shrink]
func (fsa *FSAdmin) ResizeSubVolumeGroup(volume, group string,
	newSize QuotaSize, noShrink bool) (*SubVolumeGroupResizeResult, error) {

	f := &subVolumeGroupResizeFields{
		Prefix:    "fs subvolumegroup resize",
		Format:    "json",
		VolName:   volume,
		GroupName: group,
		NewSize:   newSize.resizeValue(),
		NoShrink:  noShrink,
	}

	var rawData []map[string]interface{}
	var result SubVolumeGroupResizeResult
	res := fsa.marshalMgrCommand(f)
	if err := res.NoStatus().Unmarshal(&rawData).End(); err != nil {
		return nil, err
	}
	for _, item := range rawData {
		for key, value := range item {
			switch key {
			case "bytes_used":
				result.BytesUsed = ByteCount(uint64(value.(float64))) // 将 float64 转换为 ByteCount
			case "bytes_quota":
				result.BytesQuota = ByteCount(uint64(value.(float64)))
			case "bytes_pcent":
				result.BytesPercent = value.(string)
			}
		}
	}

	return &result, nil

}
