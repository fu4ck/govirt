package volume

import libvirtgo "github.com/libvirt/libvirt-go"

type Storage struct {
	Storage *libvirtgo.StoragePool
}

func NewStoragePool() *Storage {
	return &Storage{
		Storage: &libvirtgo.StoragePool{

		},
	}
}

func (s *Storage) CreateStorage(flags libvirtgo.StoragePoolCreateFlags) error {
	err := s.Storage.Create(flags)
	if err != nil {
		return err
	}
	return nil
}


func (s *Storage) RemoveStorage(flags libvirtgo.StoragePoolDeleteFlags) error {
	err := s.Storage.Delete(flags)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetInfo() (*libvirtgo.StoragePoolInfo, error) {
	info, err := s.Storage.GetInfo()
	if err != nil {
		return nil, err
	}
	return info, nil
}