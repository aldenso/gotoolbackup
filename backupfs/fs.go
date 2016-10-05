package backupfs

import (
	"github.com/spf13/afero"
)

var sourceFs afero.Fs

// InitMemFs initializes a memory fs.
func InitMemFs() afero.Fs {
	sourceFs := &afero.MemMapFs{}
	return sourceFs
}

// InitOSFs initializes an OS fs.
func InitOSFs() afero.Fs {
	sourceFs = &afero.OsFs{}
	return sourceFs
}
