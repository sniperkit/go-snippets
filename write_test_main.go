package main

import "os"

func writeTestMain(file *os.File) (err error) {
	_, err = file.WriteString(`#define CATCH_CONFIG_MAIN
#include <catch.hpp>`)

	return
}
