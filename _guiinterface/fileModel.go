package main

import "github.com/therecipe/qt/core"

const (
	FilePath = int(core.Qt__UserRole) + 1<<iota
	FileSize
	FileSource
)

type FileModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*File                  `property:"files"`

	_ func(*File)                                          `slot:"addFile"`
	_ func(row int, filePath, fileSize, fileSource string) `slot:"editFile"`
	_ func(row int)                                        `slot:"removeFile"`
	_ func()                                               `slot:"clearFiles"`
}

type File struct {
	core.QObject

	_ string `property:"filePath"`
	_ string `property:"fileSize"`
	_ string `property:"fileSource"`
}

func init() {
	File_QRegisterMetaType() //was Person_ i changed it, if any errors...
}

func (m *FileModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		FilePath:   core.NewQByteArray2("filePath", len("filePath")),
		FileSize:   core.NewQByteArray2("fileSize", len("fileSize")),
		FileSource: core.NewQByteArray2("fileSource", len("fileSource")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddFile(m.addFile)
	m.ConnectEditFile(m.editFile)
	m.ConnectRemoveFile(m.removeFile)
	m.ConnectClearFiles(m.clearFiles)

}

func (m *FileModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.Files()) {
		return core.NewQVariant()
	}

	var f = m.Files()[index.Row()]

	switch role {
	case FilePath:
		{
			return core.NewQVariant14(f.FilePath())
		}
	case FileSize:
		{
			return core.NewQVariant14(f.FileSize())
		}
	case FileSource:
		{
			return core.NewQVariant14(f.FileSource())
		}

	default:
		{
			return core.NewQVariant()
		}
	}
}

func (m *FileModel) rowCount(parent *core.QModelIndex) int {
	return len(m.Files())
}

func (m *FileModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *FileModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *FileModel) addFile(f *File) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.Files()), len(m.Files()))
	m.SetFiles(append(m.Files(), f))
	m.EndInsertRows()
}

func (m *FileModel) editFile(row int, filePath string, fileSize string, fileSource string) {
	var p = m.Files()[row]

	if filePath != "" {
		p.SetFilePath(filePath)
	}
	if fileSize != "" {
		p.SetFileSize(fileSize)
	}
	if fileSource != "" {
		p.SetFileSource(fileSource)
	}

	var fIndex = m.Index(row, 0, core.NewQModelIndex())
	m.DataChanged(fIndex, fIndex, []int{FilePath, FileSize, FileSource})
}
func (m *FileModel) clearFiles() {
	m.BeginResetModel()
	m.SetFiles(make([]*File, 0))
	m.EndResetModel()
}
func (m *FileModel) removeFile(row int) {
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetFiles(append(m.Files()[:row], m.Files()[row+1:]...))
	m.EndRemoveRows()
}
