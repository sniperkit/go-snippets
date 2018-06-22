/*
* @Author: souravray
* @Date:   2014-10-21 03:35:32
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 22:17:36
 */

package models

import (
	"errors"
	"regexp"
)

var (
	EmailRegexp   = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	URLRootRegexp = regexp.MustCompile(`^[a-z0-9.\-]+$`)
)

var (
	ErrInvalidEmail    = errors.New(`Invalid email`)
	ErrInvalidName     = errors.New(`Invalid name`)
	ErrInvalidPassword = errors.New(`Password must be of 4 charecters and more`)
	ErrInvalidURLRoot  = errors.New(`Your board url can only contain alpha numaric charecters and - `)
	ErrNotFilled       = errors.New(`Blank email`)
)

func init() {

}
