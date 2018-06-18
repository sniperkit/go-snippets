package main

// struct ata_identify_device {
//   unsigned short words000_009[10];
//   unsigned char  serial_no[20];
//   unsigned short words020_022[3];
//   unsigned char  fw_rev[8];
//   unsigned char  model[40];
//   unsigned short words047_079[33];
//   unsigned short major_rev_num;
//   unsigned short minor_rev_num;
//   unsigned short command_set_1;
//   unsigned short command_set_2;
//   unsigned short command_set_extension;
//   unsigned short cfs_enable_1;
//   unsigned short word086;
//   unsigned short csf_default;
//   unsigned short words088_255[168];
// } ATTR_PACKED;
//
// void swapbytes(char *out, char *in, int n) {
//   int i;
//   for (i = 0; i < n; i+=2) {
//      out[i]   = in[i+1];
//      out[i+1] = in[i];
//   }
// }
//
// struct ata_identify_device *convert_ata_identify_device(char *out, char *data, int n) {
//   swapbytes(out, data, n);
//   return (struct ata_identify_device*)(out);
// }
import "C"

import (
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
)

const (
	DevSdb string = "/dev/sdb"
)

var (
	ATA_IDENTIFY_DEVICE uint8 = 0xec
)

func main() {
	log.Printf("Prepare to open device %s\n", DevSdb)
	file, err := os.Open(DevSdb)
	if err != nil {
		panic(err)
	}

	// unsigned char deviceid[4*sizeof(int)+512*sizeof(char)]
	var p [516]byte
	p[0] = ATA_IDENTIFY_DEVICE
	p[3] = 1

	_p0 := unsafe.Pointer(&p[0])

	r1, r2, errno := syscall.Syscall6(syscall.SYS_IOCTL, file.Fd(), uintptr(0x31F), uintptr(_p0), uintptr(len(p)), 0, 0)
	if errno != 0 {
		panic(errno)
	}
	spew.Dump(r1)
	spew.Dump(r2)
	spew.Dump(errno)
	spew.Dump(p[4:])

	data := p[4:]
	out := make([]byte, len(data))
	identify := C.convert_ata_identify_device(C.CString(string(out)), C.CString(string(data)), C.int(len(data)))
	spew.Dump(identify)

	<-(chan struct{})(nil)
}
