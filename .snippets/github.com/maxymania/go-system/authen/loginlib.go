/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 Checks user Credentials (/etc/shaddow). Uses "github.com/kless/osutil/user" and
 "github.com/kless/osutil/user/crypt".
 */
package authen

import "github.com/kless/osutil/user"
import "github.com/kless/osutil/user/crypt"
import _ "github.com/kless/osutil/user/crypt/apr1_crypt"
import _ "github.com/kless/osutil/user/crypt/md5_crypt"
import _ "github.com/kless/osutil/user/crypt/sha256_crypt"
import _ "github.com/kless/osutil/user/crypt/sha512_crypt"

import "strings"
import "errors"
import "fmt"

var NoSuchUser = errors.New("no such user")

func catch2error(i interface{}, e *error) {
	if i==nil { return }
	if v,ok := i.(error); ok { *e = v; return }
	*e = errors.New(fmt.Sprint(i))
}

/*
 Authenticates an user using his name and password. Returns nil, if the
 Credentials match.
 */
func AuthenticatePassword(usr string, password []byte) (err error) {
	defer func() { catch2error(recover(),&err) }()
	s,e := user.LookupShadow(usr)
	if e!=nil { return e }
	parts := strings.Split(s.String(),":")
	if len(parts)<2 { return NoSuchUser }
	return crypt.NewFromHash(parts[1]).Verify(parts[1],password)
}


