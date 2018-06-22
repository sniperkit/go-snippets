// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"os"
	"errors"
	"sync"
	"encoding/binary"
)

type BinLog struct {
	fileName      string
	file          *os.File
	writer        *bufio.Writer
	mutex			sync.Mutex
}

func openBinLog(fileName string) (*BinLog, error) {
	file, err := os.OpenFile(fileName, os.O_APPEND, 0644)
	if err != nil {
		file, err = os.Create(fileName)
		if err != nil {
			return nil, errors.New("can not create the log file.")
		}
	}
	bg := &BinLog{}
	bg.file = file
	bg.writer = bufio.NewWriter(file)
	return bg, nil
}

func (bg *BinLog) append(version uint32, data []byte) {
	bg.mutex.Lock()
	defer func() {
		bg.mutex.Unlock()
	}()
	// 4 is data length
	// 4 is version
	lenBuff := make([]byte, 8)
	binary.BigEndian.PutUint32(lenBuff[:4], uint32(len(data)))
	binary.BigEndian.PutUint32(lenBuff[4:], version)

	bg.writer.Write(lenBuff)
	bg.writer.Write(data)
	// it need to write safe,so to flush each time.
	bg.writer.Flush()
}

func (bg *BinLog) close() {
	bg.mutex.Lock()
	defer func() {
		bg.mutex.Unlock()
	}()
	bg.writer.Flush()
	bg.file.Close()
}

type BinLogReader struct {
	fileName string
	file          *os.File
	reader *bufio.Reader
	index int64
}

func openBinLogReader(fileName string) (*BinLogReader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, errors.New("can not open the log file.")
	}
	blr := &BinLogReader{}
	blr.file = file
	blr.reader = bufio.NewReader(file)
	return blr, nil
}

func (blr *BinLogReader) next() (version uint32, data []byte, err error) {
	versionBuf := make([]byte, 8)
	err = readLength(blr.reader, versionBuf, uint32(8))
	if err != nil {
		return 0, nil, err
	}
	length := binary.BigEndian.Uint32(versionBuf[:4])
	version = binary.BigEndian.Uint32(versionBuf[4:])
	data = make([]byte, length)
	err = readLength(blr.reader, data, length)
	if err != nil {
		return 0, nil, err
	}
	return version, data, nil
}

func (blr *BinLogReader) close() {
	blr.file.Close()
}

func readLength(reader *bufio.Reader, buf []byte, length uint32) (error) {
	var index uint32
	index = 0;
	times := 0;
	for {
		if times > 10 {
			return errors.New("can not read full in 10 times.")
		}
		count, err := reader.Read(buf[index:])
		if err != nil{
			return err
		}
		index = uint32(count) + index
		if index >= length {
			break
		}
		times++
	}
	return nil
}
