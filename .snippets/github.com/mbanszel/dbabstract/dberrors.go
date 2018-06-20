// Copyright (c) 2018 Jef Oliver. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dbabstract

import "errors"

var (
	// ErrDatabaseConnectTypeUnsupported is the error returned for unsupported network types
	ErrDatabaseConnectTypeUnsupported = errors.New("unsupported connection type")
	// ErrDatabaseConnectTypeMissing is the error returned for missing database connection types
	ErrDatabaseConnectTypeMissing = errors.New("no database connection type provided")
	// ErrDatabaseHostMissing is the error returned for missing database hosts
	ErrDatabaseHostMissing = errors.New("no database host provided")
	// ErrDatabaseNameMissing is the error returned for missing database names
	ErrDatabaseNameMissing = errors.New("no database name provided")
	// ErrDatabasePortMissing is the error returned for missing database ports
	ErrDatabasePortMissing = errors.New("no database port provided")
	// ErrDatabaseSocketPathMissing is the error returned for missing database socket paths
	ErrDatabaseSocketPathMissing = errors.New("no database socket path provided")
	// ErrDataPathMissing is returned when no data path was provided
	ErrDataPathMissing = errors.New("no data path provided")
	// ErrDataPathNotDir is returned when a path exists but is not a directory
	ErrDataPathNotDir = errors.New("data path exists but is not a directory")
)
