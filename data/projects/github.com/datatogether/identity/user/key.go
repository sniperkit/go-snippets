package user

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/datatogether/errors"
	"github.com/datatogether/sqlutil"
	"time"

	"golang.org/x/crypto/ssh"
)

type Key struct {
	Type     string
	Sha256   [32]byte
	Created  int64
	LastSeen int64
	Name     string
	User     *User
	Public   []byte
	private  []byte
}

func (key *Key) MarshalJSON() ([]byte, error) {
	// TODO - switch this to an override of Key.Sha256's marshal method
	return json.Marshal(map[string]interface{}{
		"type":      key.Type,
		"sha256":    fmt.Sprintf("%x", key.Sha256),
		"last_seen": key.LastSeen,
		"name":      key.Name,
		"user":      key.User,
	})
}

func keyString(pubKey ssh.PublicKey) string {
	return fmt.Sprintf("%x", sha256.Sum256(pubKey.Marshal()))
}

// UserForPublicKey finds a user based on a provided public key
func UserForPublicKey(db sqlutil.Queryable, pubKey ssh.PublicKey) (*User, error) {
	row := db.QueryRow(fmt.Sprintf("select %s from users where id = (select id from keys where sha_256= $1 and deleted = false)", userColumns()), sha256.Sum256(pubKey.Marshal()))
	u := &User{}
	err := u.UnmarshalSQL(row)
	return u, err
}

// TODO
func (u *User) Keys(db *sql.DB) ([]*Key, error) {
	rows, err := db.Query(fmt.Sprintf("select %s from keys where user_id= $1 and deleted = false", keyColumns()), u.Id)
	if err != nil {
		return nil, err
	}
	keys := []*Key{}
	for rows.Next() {
		k := &Key{}
		if err := k.UnmarshalSQL(rows); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// MakeKeyPair generates an
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
// func MakeKeyPair(pubKeyPath, privateKeyPath string) error {
// 	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
// 	if err != nil {
// 		return err
// 	}

// 	// generate and write private key as PEM
// 	privateKeyFile, err := os.Create(privateKeyPath)
// 	defer privateKeyFile.Close()
// 	if err != nil {
// 		return err
// 	}
// 	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
// 	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
// 		return err
// 	}

// 	// generate and write public key
// 	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
// }

func NewKey(name string, u *User) (*Key, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &Key{
		Sha256:  sha256.Sum256(pub.Marshal()),
		Type:    "ssh",
		User:    u,
		Name:    name,
		Created: time.Now().Unix(),
		Public:  pub.Marshal(),
		private: pem.EncodeToMemory(privateKeyPEM),
	}, nil
}

func CreateKey(db sqlutil.Execable, u *User, name string, publicKey []byte) (*Key, error) {
	key, comment, _, _, err := ssh.ParseAuthorizedKey(publicKey)
	if err != nil {
		return nil, errors.ErrInvalidKey
	}

	if name == "" && comment != "" {
		name = comment
	}

	pkBytes := key.Marshal()
	k := &Key{
		Type:     key.Type(),
		User:     u,
		Sha256:   sha256.Sum256(pkBytes),
		Created:  time.Now().Unix(),
		Name:     name,
		LastSeen: 0,
		Public:   pkBytes,
	}

	if err := k.validate(db); err != nil {
		return nil, err
	}

	if _, e := db.Exec("INSERT INTO keys VALUES ($1, $2, $3, $4, $5, $6, $7, null, false)", k.Type, k.Sha256[:], k.Created, k.LastSeen, k.Name, k.User.Id, k.Public); e != nil {
		return nil, errors.NewFmtError(500, e.Error())
	}

	return k, nil
}

func (key *Key) Read(db sqlutil.Queryable) error {
	row := db.QueryRow(fmt.Sprintf("SELECT %s FROM keys WHERE sha_256=$1 AND deleted=false", keyColumns()), key.Sha256[:])
	if err := key.UnmarshalSQL(row); err != nil {
		if err == sql.ErrNoRows {
			return errors.ErrNotFound
		} else {
			return errors.New500Error(err.Error())
		}
	}
	return nil
}

func (k *Key) Save(db sqlutil.Execable) error {
	if err := k.validate(db); err != nil {
		return err
	}

	prev := &Key{Sha256: k.Sha256}
	if err := prev.Read(db); err != nil {
		if err == errors.ErrNotFound {
			k.Created = time.Now().Unix()
			if _, e := db.Exec("INSERT INTO keys VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false)", k.Type, k.Sha256[:], k.Created, k.LastSeen, k.Name, k.User.Id, k.Public, k.private); e != nil {
				return errors.NewFmtError(500, e.Error())
			}
		} else {
			return err
		}
	} else {
		_, err := db.Exec("UPDATE keys SET last_seen = $2 WHERE sha_256 = $1 AND deleted = false", k.Sha256[:], k.LastSeen)
		return err
	}

	return nil
}

// "delete" a user_key
func (key *Key) Delete(db sqlutil.Execable) error {
	_, err := db.Exec("UPDATE keys SET deleted=true WHERE sha_256 = $1", key.Sha256[:])
	return err
}

func (key *Key) validate(db sqlutil.Queryable) error {
	if key.User == nil {
		return errors.ErrUserRequired
	}

	if key.Name == "" {
		return errors.ErrNameRequired
	}

	// TODO - implement UserExists
	// if exists, err := UserExists(db, key.User.Id); err != nil {
	// 	return err
	// } else if !exists {
	// 	return ErrInvalidUser
	// }

	return nil
}

func keyColumns() string {
	return "type, sha_256, created, last_seen, name, user_id, public_bytes, private_bytes"
}

// turn an sql row from the user table into a user struct pointer
func (key *Key) UnmarshalSQL(row sqlutil.Scannable) error {
	var (
		keyType, name, userId             string
		created, lastSeen                 int64
		keySha, publicBytes, privateBytes []byte
		keySha256                         = [32]byte{}
	)

	// "type, key_sha, created, last_seen, name, user_id, public_bytes, private_bytes"
	if err := row.Scan(&keyType, &keySha, &created, &lastSeen, &name, &userId, &publicBytes, &privateBytes); err != nil {
		return err
	}

	for i, b := range keySha {
		keySha256[i] = b
	}

	*key = Key{
		Type:     keyType,
		Sha256:   keySha256,
		Created:  created,
		LastSeen: lastSeen,
		Name:     name,
		User:     NewUserFromString(userId),
		Public:   publicBytes,
		private:  privateBytes,
	}

	return nil
}
