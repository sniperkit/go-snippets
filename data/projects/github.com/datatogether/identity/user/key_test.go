package user

import (
	"fmt"
	// "testing"
)

func CompareKeys(a, b *Key) error {
	if a.Sha256 != b.Sha256 {
		return fmt.Errorf("mismatch Sha256: %x != %x", a.Sha256, b.Sha256)
	}

	return nil
}

func generateTestPublicKey() []byte {
	// TODO - learn how to do this properly
	// private, err := rsa.GenerateKey(rand.Reader, 1024)
	// if err != nil {
	// 	panic(err)
	// }
	// private.Public()
	// pub, err := ssh.NewPublicKey(private.Public())
	// if err != nil {
	// 	panic(err)
	// }
	// return pub.Marshal()
	return []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDTpxZFUdgKDSdQPZR3DeVmjgcKlIfzCoSu3tyaImoknfCCJX/ZMeuxbq2+pSHwNjfmwiJsT8zwRfGvzOmXqVtiQGZN1F7AZK53JlXrp5YbFOfhM6eDxP1vN3/F1PVKip78iB63MUUf2ySYOM5gCgm2WOYnydSBEMnK2AQOUEeM86+vOIOq6bfO4sXcS+8CyRiV/RYjHkoNU3pdWIjFSHKkI4tY05wPFJpT636G7qa1V6MwN3cM0qrxuZpQ4Cc+MeOkbAHJB2oEvN10JOFEmE4pIBSy33ZWRnhMy/5Yzd85Suenfe3D9A3fu3lZXdm/1li7UdzHm81WASX8cP/VlH3V brendan@Brendans-MacBook-Pro.local")
}

// func TestKey(t *testing.T) {
// 	cases := []struct {
// 		user                                *User
// 		name                                string
// 		key                                 []byte
// 		createErr, readErr, saveErr, delErr error
// 	}{
// 		{nil, "", generateTestPublicKey(), ErrUserRequired, nil, nil, nil},
// 		// TODO - name is now derived from key comment.
// 		// {TestData.Users.Brendan, "", generateTestPublicKey(), ErrNameRequired, nil, nil, nil},
// 		// {TestData.Users.Brendan, "", nil, ErrInvalidKey, nil, nil, nil},
// 		// {TestData.Users.Brendan, "test key", generateTestPublicKey(), nil, nil, nil, nil},
// 	}

// 	for i, c := range cases {
// 		key, got := CreateKey(testDB, c.user, c.name, c.key)
// 		if got != c.createErr {
// 			t.Errorf("case %d create key error mismatch. expected: '%s', got: '%s'", i, c.createErr, got)
// 		}
// 		if c.createErr != nil || got != c.createErr {
// 			continue
// 		}

// 		r := &Key{Sha256: key.Sha256}
// 		if got = r.Read(testDB); got != c.readErr {
// 			t.Errorf("case %d read error mismatch. expected: '%s', got: '%s'", i, c.readErr, got)
// 		}
// 		if got = CompareKeys(key, r); got != nil {
// 			t.Errorf("case %d read mismatch: '%s'", i, got)
// 		}

// 		if got = r.Save(testDB); got != c.saveErr {
// 			t.Errorf("case %d save error mismatch. expected: '%s', got: '%s'", i, c.saveErr, got)
// 		}

// 		if got = r.Delete(testDB); got != c.delErr {
// 			t.Errorf("case %d delete error mismatch. expected: '%s', got: '%s'", i, c.delErr, got)
// 		}
// 	}
// }
