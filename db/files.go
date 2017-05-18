package db

import "database/sql"

// File represents a file ID to name link.
type File struct {
	ID        string
	Name      string
	Namespace string
	Extension string
}

const filesSchema = `
	id CHAR(32) NOT NULL,
	name VARCHAR(255) NOT NULL,
	namespace VARCHAR(255) NOT NULL,
	extension VARCHAR(255) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY (name, namespace)
`

// GetFileByID gets a file by its storage ID.
func GetFileByID(id string) *File {
	row := db.QueryRow(`SELECT * FROM files WHERE id=?`, id)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

// GetFileByPath gets a file by its namespace, name and extension.
func GetFileByPath(name, namespace string) *File {
	row := db.QueryRow(`SELECT * FROM files WHERE name=? AND namespace=?`, name, namespace)
	if row != nil {
		return scanFile(row)
	}
	return nil
}

func scanFile(row *sql.Row) *File {
	var id, name, namespace, extension string
	row.Scan(&id, &name, &namespace, &extension)
	return &File{ID: id, Name: name, Namespace: namespace, Extension: extension}
}

// GetPermissions returns the permissions to this file.
func (file *File) GetPermissions() []*Permission {
	results, err := db.Query(`SELECT * FROM permissions WHERE file=?`, file.ID)
	if err != nil {
		return []*Permission{}
	}
	return scanPermissions(results)
}
