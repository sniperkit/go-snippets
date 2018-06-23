# mapdecor
Decorator of map[string]interface{} registry and predefined decorators

[![Build Status](https://travis-ci.org/linkosmos/mapdecor.svg?branch=master)](https://travis-ci.org/linkosmos/mapdecor)
[![Coverage Status](https://coveralls.io/repos/github/linkosmos/mapdecor/badge.svg?branch=master)](https://coveralls.io/github/linkosmos/mapdecor?branch=master)
[![GoDoc](http://godoc.org/github.com/linkosmos/mapdecor?status.svg)](http://godoc.org/github.com/linkosmos/mapdecor)
[![Go Report Card](http://goreportcard.com/badge/linkosmos/mapdecor)](http://goreportcard.com/report/linkosmos/mapdecor)
[![BSD License](http://img.shields.io/badge/license-BSD-blue.svg)](http://opensource.org/licenses/BSD-3-Clause)

### Example

```go
    input := map[string]interface{}{
      "key1": nil,
      "key2": "with",
      "val1": "1",
      "val2": "2",
      "val3": "3",
      "val4": "4",
    }

    decorFunc := func(input map[string]interface{}) (output map[string]interface{}) {
      partitonFunc := func(s string, i interface{}) bool {
        return strings.Contains(s, "val")
      }

      // For first (inputPartitioned[0]) partition we get key-values containing 'val' in key
      // For second (inputPartitioned[1]) partition what is left
      inputPartitioned := mapop.Partition(partitonFunc, input)

      // Assigning values key to first partition
      inputPartitioned[1]["values"] = inputPartitioned[0]

      return inputPartitioned[1]
    }


    got := Decorate(input, decorFunc)

  // Got
  // map[string]interface{}{
  //   "key1": nil,
  //   "key2": "with",
  //   "values": map[string]interface{}{
  //     "val1": "1",
  //     "val2": "2",
  //     "val3": "3",
  //     "val4": "4",
  //   }
  // }

```

### License

Copyright (c) 2015, linkosmos
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of mapdecor nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
