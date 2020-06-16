package storage

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"

	"sb.im/gosd/model"
)

func (s *Storage) CreateBlob(blob *model.Blob) (err error) {
	query := `
	INSERT INTO blobs
	(filename, content, checksum, create_at, update_at)
	VALUES
	($1, $2, $3, now(), now())
	RETURNING
	id, filename
	`

	data, err := ioutil.ReadAll(blob.Reader)
	if err != nil {
		return fmt.Errorf(`store: unable to read io: %v`, err)
	}

	err = s.db.QueryRow(query, blob.FileName, data, fmt.Sprintf("%x", sha1.Sum(data))).Scan(
		&blob.ID,
		&blob.FileName,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to create blob: %v`, err)
	}

	return nil
}
