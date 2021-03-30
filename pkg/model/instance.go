package model

type DiskList struct {
	Disk string
}

type Instance struct {
	HostName   string
	Id         string
	LocalMount []Mount
	HostServer string
	VNCUrl     string
}

type Mount struct {
	SourcePoint string
}


func NewInstance() *Instance {
	return &Instance{
		HostName:   "1111",
		Id:         "2222",
		LocalMount: nil,
		HostServer: "",
		VNCUrl:     "",
	}
}