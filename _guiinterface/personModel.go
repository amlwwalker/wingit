package main

import "github.com/therecipe/qt/core"

const (
	FirstName = int(core.Qt__UserRole) + 1<<iota
	LastName
	Email
	FileCount
)

type PersonModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Person                `property:"people"`

	_ func(*Person)                                    `slot:"addPerson"`
	_ func(row int, firstName, lastName, email, fileCount string) `slot:"editPerson"`
	_ func(row int)                                    `slot:"removePerson"`
	_ func()                                    		`slot:"clearPeople"`
}

type Person struct {
	core.QObject

	_ string `property:"firstName"`
	_ string `property:"lastName"`
	_ string `property:"email"`
	_ string `property:"fileCount"`
}

func init() {
	Person_QRegisterMetaType()
}

func (m *PersonModel) init() {
	m.SetRoles(map[int]*core.QByteArray{
		FirstName: core.NewQByteArray2("firstName", len("firstName")),
		LastName:  core.NewQByteArray2("lastName", len("lastName")),
		Email:     core.NewQByteArray2("email", len("email")),
		FileCount:     core.NewQByteArray2("fileCount", len("fileCount")),
	})

	m.ConnectData(m.data)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectRoleNames(m.roleNames)

	m.ConnectAddPerson(m.addPerson)
	m.ConnectEditPerson(m.editPerson)
	m.ConnectRemovePerson(m.removePerson)
	m.ConnectClearPeople(m.clearPeople)
}

func (m *PersonModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	if index.Row() >= len(m.People()) {
		return core.NewQVariant()
	}

	var p = m.People()[index.Row()]

	switch role {
	case FirstName:
		{
			return core.NewQVariant14(p.FirstName())
		}
	case LastName:
		{
			return core.NewQVariant14(p.LastName())
		}
	case Email:
		{
			return core.NewQVariant14(p.Email())
		}
	case FileCount:
		{
			return core.NewQVariant14(p.FileCount())
		}
	default:
		{
			return core.NewQVariant()
		}
	}
}

func (m *PersonModel) rowCount(parent *core.QModelIndex) int {
	return len(m.People())
}

func (m *PersonModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (m *PersonModel) roleNames() map[int]*core.QByteArray {
	return m.Roles()
}

func (m *PersonModel) addPerson(p *Person) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.People()), len(m.People()))
	m.SetPeople(append(m.People(), p))
	m.EndInsertRows()
}

func (m *PersonModel) editPerson(row int, firstName, lastName, email, fileCount string) {
	var p = m.People()[row]

	if firstName != "" {
		p.SetFirstName(firstName)
	}
	if lastName != "" {
		p.SetLastName(lastName)
	}
	if email != "" {
		p.SetEmail(email)
	}
	if fileCount != "" {
		p.SetFileCount(fileCount)
	}
	var pIndex = m.Index(row, 0, core.NewQModelIndex())
	m.DataChanged(pIndex, pIndex, []int{FirstName, LastName, Email, FileCount})
}

func (m *PersonModel) clearPeople() {
	m.BeginResetModel()
	m.SetPeople(make([]*Person, 0))
	m.EndResetModel()
}

func (m *PersonModel) removePerson(row int) {
	m.BeginRemoveRows(core.NewQModelIndex(), row, row)
	m.SetPeople(append(m.People()[:row], m.People()[row+1:]...))
	m.EndRemoveRows()
}
