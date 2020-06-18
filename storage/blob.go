package storage

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
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

func (s *Storage) BlobByID(id int64) (*model.Blob, error) {
	query := `
		SELECT
			id,
			filename,
			content
		FROM
			blobs
		WHERE
			id = $1
	`

	return s.fetchBlob(query, id)
}

func (s *Storage) fetchBlob(query string, args ...interface{}) (*model.Blob, error) {
	blob := &model.Blob{}

	var content []byte

	err := s.db.QueryRow(query, args...).Scan(
		&blob.ID,
		&blob.FileName,
		&content,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch blob: %v`, err)
	}

	blob.Reader = bytes.NewReader(content)

	return blob, nil
}
