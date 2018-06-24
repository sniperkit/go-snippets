package core

import (
	"github.com/datatogether/sql_datastore"
	"testing"
)

func metadataEqual(a, b *Metadata) bool {
	if a.Hash != "" && b.Hash != "" {
		return a.Hash == b.Hash
	}

	return a.Timestamp.Equal(b.Timestamp) && a.Subject == b.Subject && a.KeyId == b.KeyId && a.Prev == b.Prev
}

func TestMetadata(t *testing.T) {
	defer resetTestData(appDB, "metadata")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Url{}); err != nil {
		t.Error(err.Error())
		return
	}

	keyId := "test_key_id"
	subject := "test_subject"

	m, err := NextMetadata(appDB, keyId, subject)
	if err != nil {
		t.Error(err.Error())
		return
	}

	m.Meta = map[string]interface{}{
		"key": "value",
	}

	if err := m.Write(store); err != nil {
		t.Error(err.Error())
		return
	}

	b := &Metadata{
		Timestamp: m.Timestamp,
		KeyId:     keyId,
		Subject:   subject,
		Meta: map[string]interface{}{
			"key": "value",
		},
	}

	if !metadataEqual(m, b) {
		t.Errorf("metdata mismach: %s != %s", m, b)
	}

	c, err := LatestMetadata(appDB, keyId, subject)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !metadataEqual(m, c) {
		t.Errorf("metadta mismatch: %s != %s ", m, c)
	}
}
