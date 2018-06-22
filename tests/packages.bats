#!/usr/bin/env bats

# TODO: Auto generate this? Do so in topological order to pinpoint breaking libraries.

# TODO: Revisit using `go test`.

## Libraries (no dependencies)

@test "Install zlib" {
	gogurt zlib
}

@test "Install bzip2" {
	gogurt bzip2
}

@test "Install gmp" {
	gogurt gmp
}

@test "Install mpfr" {
	gogurt mpfr
}

@test "Install mpc" {
	gogurt mpc
}

@test "Install OpenSSL" {
	gogurt openssl
}

# "Small" utilities (have one or two dependencies)
@test "Install pigz" {
	gogurt pigz
}

@test "Install pbzip2" {
	gogurt pbzip2
}

# Takes too long to install.
# @test "Install gcc" {
# 	gogurt gcc
# }
