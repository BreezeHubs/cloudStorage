package meta

type FileMeta struct {
	FileName string
	FileSize int64
	FileSha1 string
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta updates file meta to map
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta gets file meta from map
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaList gets file meta list from map
func GetFileMetas() map[string]FileMeta {
	return fileMetas
}

// DeleteFileMeta deletes file meta from map
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
