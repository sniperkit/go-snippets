package bfd

/*
#include <bfd.h>
#include <stdio.h>
#include <stdlib.h>
#include "gobfd.h"

#cgo LDFLAGS: -lbfd
*/
import "C"

import (
	"math"
	"os"
	"sync"
	"unsafe"
)

func init() {
	C.bfd_init()
}

type (
	File         C.bfd
	Target       C.bfd_target
	VMA          C.bfd_vma
	Section      C.asection
	PluginFormat C.enum_bfd_plugin_format
	BuildID      C.struct_bfd_build_id
	Flagword     C.flagword
	Size         C.bfd_size_type
	Symbol       C.asymbol
	SymbolTable  struct {
		syms  unsafe.Pointer
		size  int64
		count int64
	}
	sectionHandler struct{ handler func(*File, *Section) }
)

func (c *File) Filename() string   { return C.GoString(c.filename) }
func (c *File) Sections() *Section { return (*Section)(c.sections) }
func (c *File) Xvec() *Target      { return (*Target)(c.xvec) }
func (c *File) Flags() int         { return int(C.getFileFlags((*C.struct_bfd)(c))) }
func (c *File) SetFlags(flags int) { C.setFileFlags((*C.struct_bfd)(c), C.int(flags)) }

func (s *Section) Name() string    { return C.GoString(s.name) }
func (s *Section) LMA() VMA        { return VMA(s.lma) }
func (s *Section) VMA() VMA        { return VMA(s.vma) }
func (s *Section) Size() int64     { return int64(s.size) }
func (s *Section) Next() *Section  { return (*Section)(s.next) }
func (s *Section) Prev() *Section  { return (*Section)(s.prev) }
func (s *Section) Flags() Flagword { return Flagword(s.flags) }

func (t *Target) Name() string { return C.GoString(t.name) }

func (s *SymbolTable) Size() int64 { return s.count }
func (s *SymbolTable) Free()       { C.free(s.syms) }

func xtrue(cond C.bfd_boolean) error {
	if cond != 0 {
		return nil
	}
	return Error(C.bfd_get_error())
}

func stringList(list **C.char) []string {
	if list == nil {
		return nil
	}
	xlist := (*[math.MaxInt32]*C.char)(unsafe.Pointer(list))
	var str []string
	for i := 0; xlist[i] != nil; i++ {
		str = append(str, C.GoString(xlist[i]))
	}
	C.free(unsafe.Pointer(list))
	return str
}

type sectionMappers struct {
	sync.Mutex
	count int
	funcs []func(*File, *Section)
}

func (s *sectionMappers) Acquire(f func(*File, *Section)) unsafe.Pointer {
	s.Lock()
	defer s.Unlock()
	s.funcs = append(s.funcs, f)
	s.count++
	return unsafe.Pointer(uintptr(len(s.funcs)) - 1)
}

func (s *sectionMappers) Release() {
	s.Lock()
	defer s.Unlock()
	if s.count--; s.count == 0 {
		s.funcs = s.funcs[:0]
	}
}

var (
	sm sectionMappers
)

//export goMapOverSections
func goMapOverSections(abfd *C.bfd, section *C.asection, data unsafe.Pointer) {
	sm.Lock()
	f := sm.funcs[int(uintptr(data))]
	sm.Unlock()
	f((*File)(abfd), (*Section)(section))
}

func MapOverSections(abfd *File, f func(*File, *Section)) {
	C.bfd_map_over_sections((*C.bfd)(abfd), (*[0]byte)(C.mapOverSections), sm.Acquire(f))
	sm.Release()
}

func TargetList() []string {
	return stringList(C.bfd_target_list())
}

func SetDefaultTarget(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return xtrue(C.bfd_set_default_target(cname))
}

func FindTarget(target string, abfd *File) *Target {
	ctarget := C.CString(target)
	defer C.free(unsafe.Pointer(ctarget))
	return (*Target)(C.bfd_find_target(ctarget, (*C.bfd)(abfd)))
}

func Openr(name, target string) (*File, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var ctarget *C.char
	if target != "" {
		ctarget = C.CString(target)
		defer C.free(unsafe.Pointer(ctarget))
	}
	bfd := C.bfd_openr(cname, nil)
	return (*File)(bfd), pathError(name)
}

func Openw(name, target string) (*File, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var ctarget *C.char
	if target != "" {
		ctarget = C.CString(target)
		defer C.free(unsafe.Pointer(ctarget))
	}
	bfd := C.bfd_openw(cname, ctarget)
	return (*File)(bfd), pathError(name)
}

func Create(name string, tmpl *File) (*File, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	bfd := C.bfd_create(cname, (*C.struct_bfd)(tmpl))
	return (*File)(bfd), GetError()
}

func MakeWritable(abfd *File) error {
	return xtrue(C.bfd_make_writable((*C.bfd)(abfd)))
}

func MakeReadable(abfd *File) error {
	return xtrue(C.bfd_make_readable((*C.bfd)(abfd)))
}

func Close(abfd *File) {
	C.bfd_close((*C.bfd)(abfd))
}

func CloseAllDone(abfd *File) {
	C.bfd_close_all_done((*C.bfd)(abfd))
}

func GetArchSize(abfd *File) int {
	return int(C.bfd_get_arch_size((*C.bfd)(abfd)))
}

func GetSignExtendVMA(abfd *File) int {
	return int(C.bfd_get_sign_extend_vma((*C.bfd)(abfd)))
}

func GetSize(abfd *File) int64 {
	return int64(C.bfd_get_size((*C.bfd)(abfd)))
}

func GetMtime(abfd *File) int64 {
	return int64(C.bfd_get_mtime((*C.bfd)(abfd)))
}

func SetStartAddress(abfd *File, vma VMA) error {
	return xtrue(C.bfd_set_start_address((*C.bfd)(abfd), (C.bfd_vma)(vma)))
}

func GetStartAddress(abfd *File) VMA {
	return VMA(C.getStartAddress((*C.bfd)(abfd)))
}

func GetGPSize(abfd *File) uint {
	return uint(C.bfd_get_gp_size((*C.bfd)(abfd)))
}

func SetGPSize(abfd *File, size uint) {
	C.bfd_set_gp_size((*C.bfd)(abfd), C.uint(size))
}

func SetSectionSize(abfd *File, sec *Section, size Size) {
	C.bfd_set_section_size((*C.bfd)(abfd), (*C.asection)(sec), C.bfd_size_type(size))
}

func InitSectionDecompressStatus(abfd *File, section *Section) error {
	return xtrue(C.bfd_init_section_decompress_status((*C.bfd)(abfd), (*C.asection)(section)))
}

func InitSectionCompressStatus(abfd *File, section *Section) error {
	return xtrue(C.bfd_init_section_compress_status((*C.bfd)(abfd), (*C.asection)(section)))
}

func CheckFormat(abfd *File, format Format) error {
	return xtrue(C.bfd_check_format((*C.bfd)(abfd), C.bfd_format(format)))
}

func CheckFormatMatches(abfd *File, format Format) ([]string, error) {
	var matches **C.char
	err := xtrue(C.bfd_check_format_matches((*C.bfd)(abfd), C.bfd_format(format), &matches))
	if err != nil {
		return nil, err
	}
	return stringList(matches), nil
}

func ReadSymbolTable(abfd *File, dynamic bool) (*SymbolTable, error) {
	var size C.uint
	var count C.long
	var cdynamic C.bfd_boolean
	if dynamic {
		cdynamic = 1
	}
	syms := C.readSymbolTable((*C.bfd)(abfd), cdynamic, &size, &count)
	if count < 0 {
		return nil, GetError()
	}
	return &SymbolTable{
		syms:  unsafe.Pointer(syms),
		size:  int64(size),
		count: int64(count),
	}, nil
}

func GetSymtabUpperBound(abfd *File) int64 {
	return int64(C.getSymtabUpperBound((*C.bfd)(abfd)))
}

func GetDynamicSymtabUpperBound(abfd *File) int64 {
	return int64(C.getDynamicSymtabUpperBound((*C.bfd)(abfd)))
}

func AllocSymbolTable(size int64) *SymbolTable {
	return &SymbolTable{syms: C.malloc(C.size_t(size)), size: size}
}

func CanonicalizeSymtab(abfd *File, table *SymbolTable) (int64, error) {
	table.count = int64(C.canonicalizeSymtab((*C.struct_bfd)(abfd), (**C.struct_bfd_symbol)(table.syms)))
	if table.count < 0 {
		return 0, GetError()
	}
	return table.count, nil
}

func CanonicalizeDynamicSymtab(abfd *File, table *SymbolTable) (int64, error) {
	table.count = int64(C.canonicalizeDynamicSymtab((*C.struct_bfd)(abfd), (**C.struct_bfd_symbol)(table.syms)))
	if table.count < 0 {
		return 0, GetError()
	}
	return table.count, nil
}

func FindNearestLineDiscriminator(abfd *File, section *Section, table *SymbolTable, addr VMA) (found bool, filename, function string, line, discriminator int64) {
	var cfilename, cfunction *C.char
	var cline, cdiscriminator C.uint
	cfound := C.findNearestLineDiscriminator((*C.struct_bfd)(abfd), (*C.struct_bfd_section)(section), (**C.struct_bfd_symbol)(table.syms), C.bfd_vma(addr), &cfilename, &cfunction, &cline, &cdiscriminator)
	return cfound != 0, C.GoString(cfilename), C.GoString(cfunction), int64(cline), int64(discriminator)
}

func GetSectionByName(abfd *File, name string) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_get_section_by_name((*C.bfd)(abfd), cname))
}

func GetNextSectionByName(abfd *File, section *Section) *Section {
	return (*Section)(C.bfd_get_next_section_by_name((*C.bfd)(abfd), (*C.asection)(section)))
}

func GetLinkerSection(abfd *File, name string) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_get_linker_section((*C.bfd)(abfd), cname))
}

func GetUniqueSectionName(abfd *File, tmpl string) (string, int) {
	ctmpl := C.CString(tmpl)
	defer C.free(unsafe.Pointer(ctmpl))
	var ccount C.int
	cstr := C.bfd_get_unique_section_name((*C.bfd)(abfd), ctmpl, &ccount)
	return C.GoString(cstr), int(ccount)
}

func MakeSectionAnywayWithFlags(abfd *File, name string, flags Flagword) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_make_section_anyway_with_flags((*C.bfd)(abfd), cname, C.flagword(flags)))
}

func MakeSectionAnyway(abfd *File, name string, flags Flagword) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_make_section_anyway((*C.bfd)(abfd), cname))
}

func MakeSectionWithFlags(abfd *File, name string, flags Flagword) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_make_section_with_flags((*C.bfd)(abfd), cname, C.flagword(flags)))
}

func MakeSection(abfd *File, name string) *Section {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return (*Section)(C.bfd_make_section((*C.bfd)(abfd), cname))
}

func GetNextSectionID() int {
	return int(C.bfd_get_next_section_id())
}

func RenameSection(abfd *File, section *Section, name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.bfd_rename_section((*C.bfd)(abfd), (*C.asection)(section), cname)
}

func GenericDiscardGroup(abfd *File, group *Section) bool {
	return C.bfd_generic_discard_group((*C.bfd)(abfd), (*C.asection)(group)) != 0
}

func ScanVMA(str string, base int) (VMA, string) {
	var cend *C.char
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	vma := C.bfd_scan_vma(cstr, &cend, C.int(base))
	var xstr string
	if cend != nil {
		xstr = C.GoString(cend)
	}
	return VMA(vma), xstr
}

func SprintfVMA(abfd *File, vma VMA) string {
	var buf [80]C.char
	C.bfd_sprintf_vma((*C.bfd)(abfd), &buf[0], C.bfd_vma(vma))
	return C.GoString(&buf[0])
}

func PrintfVMA(abfd *File, vma VMA) {
	C.printfVMA((*C.bfd)(abfd), C.bfd_vma(vma))
}

func GetSectionSize(section *Section) Size {
	return Size(C.getSectionSize((*C.asection)(section)))
}

func GetSectionVMA(abfd *File, section *Section) VMA {
	return VMA(C.getSectionVMA((*C.bfd)(abfd), (*C.asection)(section)))
}

func Demangle(abfd *File, str string, options int) string {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	xstr := C.bfd_demangle((*C.struct_bfd)(abfd), cstr, C.int(options))
	defer C.free(unsafe.Pointer(xstr))
	return C.GoString(xstr)
}

func FindInlinerInfo(abfd *File) (found bool, filename, function string, line int64) {
	var cfilename, cfunction *C.char
	var cline C.uint
	cfound := C.findInlinerInfo((*C.struct_bfd)(abfd), &cfilename, &cfunction, &cline)
	return cfound != 0, C.GoString(cfilename), C.GoString(cfunction), int64(cline)
}

type Direction C.enum_bfd_direction

const (
	NoDirection    Direction = C.no_direction
	ReadDirection  Direction = C.read_direction
	WriteDirection Direction = C.write_direction
	BothDirection  Direction = C.both_direction
)

type Error C.bfd_error_type

const (
	ErrSystemCall               Error = C.bfd_error_system_call
	ErrInvalidTarget            Error = C.bfd_error_invalid_target
	ErrWrongFormat              Error = C.bfd_error_wrong_format
	ErrWrongObjectFormat        Error = C.bfd_error_wrong_object_format
	ErrInvalidOperation         Error = C.bfd_error_invalid_operation
	ErrNoMemory                 Error = C.bfd_error_no_memory
	ErrNoSymbols                Error = C.bfd_error_no_symbols
	ErrNoArmap                  Error = C.bfd_error_no_armap
	ErrNoMoreArchivedFiles      Error = C.bfd_error_no_more_archived_files
	ErrMalformedArchive         Error = C.bfd_error_malformed_archive
	ErrMissingDSO               Error = C.bfd_error_missing_dso
	ErrFileNotRecognized        Error = C.bfd_error_file_not_recognized
	ErrFileAmbguouslyRecognized Error = C.bfd_error_file_ambiguously_recognized
	ErrNoContents               Error = C.bfd_error_no_contents
	ErrNonrepresentableSection  Error = C.bfd_error_nonrepresentable_section
	ErrNoDebugSection           Error = C.bfd_error_no_debug_section
	ErrBadValue                 Error = C.bfd_error_bad_value
	ErrFileTruncated            Error = C.bfd_error_file_truncated
	ErrFileTooBig               Error = C.bfd_error_file_too_big
	ErrOnInput                  Error = C.bfd_error_on_input
	ErrInvalidCode              Error = C.bfd_error_invalid_error_code
)

func (e Error) Error() string {
	return C.GoString(C.bfd_errmsg(C.bfd_error_type(e)))
}

func GetError() error {
	err := C.bfd_get_error()
	if err == 0 {
		return nil
	}
	return Error(err)
}

func pathError(name string) error {
	err := GetError()
	if err == nil {
		return nil
	}
	return &os.PathError{Op: "open", Path: name, Err: err}
}

type Flavor C.enum_bfd_flavour

func GetFlavor(abfd *File) Flavor {
	return Flavor(C.getFlavor((*C.bfd)(abfd)))
}

const (
	TargetUnknownFlavor  Flavor = C.bfd_target_unknown_flavour
	TargetAoutFlavor     Flavor = C.bfd_target_aout_flavour
	TargetCoffFlavor     Flavor = C.bfd_target_coff_flavour
	TargetEcoffFlavor    Flavor = C.bfd_target_ecoff_flavour
	TargetXcoffFlavor    Flavor = C.bfd_target_xcoff_flavour
	TargetElfFlavor      Flavor = C.bfd_target_elf_flavour
	TargetIEEEFlavor     Flavor = C.bfd_target_ieee_flavour
	TargetNlmFlavor      Flavor = C.bfd_target_nlm_flavour
	TargetOasysFlavor    Flavor = C.bfd_target_oasys_flavour
	TargetTekhexFlavor   Flavor = C.bfd_target_tekhex_flavour
	TargetSrecFlavor     Flavor = C.bfd_target_srec_flavour
	TargetVerilogFlavor  Flavor = C.bfd_target_verilog_flavour
	TargetIhexFlavor     Flavor = C.bfd_target_ihex_flavour
	TargetSomFlavor      Flavor = C.bfd_target_som_flavour
	TargetOs9kFlavor     Flavor = C.bfd_target_os9k_flavour
	TargetVersadosFlavor Flavor = C.bfd_target_versados_flavour
	TargetMsdosFlavor    Flavor = C.bfd_target_msdos_flavour
	TargetOvaxFlavor     Flavor = C.bfd_target_ovax_flavour
	TargetEvaxFlavor     Flavor = C.bfd_target_evax_flavour
	TargetMmoFlavor      Flavor = C.bfd_target_mmo_flavour
	TargetMachoFlavor    Flavor = C.bfd_target_mach_o_flavour
	TargetPefFlavor      Flavor = C.bfd_target_pef_flavour
	TargetPefXlibFlavor  Flavor = C.bfd_target_pef_xlib_flavour
	TargetSymFlavor      Flavor = C.bfd_target_sym_flavour
)

type Endian C.enum_bfd_endian

const (
	ENDIAN_BIG     Endian = C.BFD_ENDIAN_BIG
	ENDIAN_LITTLE  Endian = C.BFD_ENDIAN_LITTLE
	ENDIAN_UNKNOWN Endian = C.BFD_ENDIAN_UNKNOWN
)

type Format C.enum_bfd_format

const (
	Unknown Format = C.bfd_unknown
	Object  Format = C.bfd_object
	Archive Format = C.bfd_archive
	Core    Format = C.bfd_core
	TypeEnd Format = C.bfd_type_end
)

const (
	SEC_NO_FLAGS                      = C.SEC_NO_FLAGS
	SEC_ALLOC                         = C.SEC_ALLOC
	SEC_LOAD                          = C.SEC_LOAD
	SEC_RELOC                         = C.SEC_RELOC
	SEC_READONLY                      = C.SEC_READONLY
	SEC_CODE                          = C.SEC_CODE
	SEC_DATA                          = C.SEC_DATA
	SEC_ROM                           = C.SEC_ROM
	SEC_CONSTRUCTOR                   = C.SEC_CONSTRUCTOR
	SEC_HAS_CONTENTS                  = C.SEC_HAS_CONTENTS
	SEC_NEVER_LOAD                    = C.SEC_NEVER_LOAD
	SEC_THREAD_LOCAL                  = C.SEC_THREAD_LOCAL
	SEC_HAS_GOT_REF                   = C.SEC_HAS_GOT_REF
	SEC_IS_COMMON                     = C.SEC_IS_COMMON
	SEC_DEBUGGING                     = C.SEC_DEBUGGING
	SEC_IN_MEMORY                     = C.SEC_IN_MEMORY
	SEC_EXCLUDE                       = C.SEC_EXCLUDE
	SEC_SORT_ENTRIES                  = C.SEC_SORT_ENTRIES
	SEC_LINK_ONCE                     = C.SEC_LINK_ONCE
	SEC_LINK_DUPLICATES               = C.SEC_LINK_DUPLICATES
	SEC_LINK_DUPLICATES_DISCARD       = C.SEC_LINK_DUPLICATES_DISCARD
	SEC_LINK_DUPLICATES_ONE_ONLY      = C.SEC_LINK_DUPLICATES_ONE_ONLY
	SEC_LINK_DUPLICATES_SAME_SIZE     = C.SEC_LINK_DUPLICATES_SAME_SIZE
	SEC_LINK_DUPLICATES_SAME_CONTENTS = C.SEC_LINK_DUPLICATES_SAME_CONTENTS
	SEC_LINKER_CREATED                = C.SEC_LINKER_CREATED
	SEC_KEEP                          = C.SEC_KEEP
	SEC_SMALL_DATA                    = C.SEC_SMALL_DATA
	SEC_MERGE                         = C.SEC_MERGE
	SEC_STRINGS                       = C.SEC_STRINGS
	SEC_GROUP                         = C.SEC_GROUP
	SEC_COFF_SHARED_LIBRARY           = C.SEC_COFF_SHARED_LIBRARY
	SEC_ELF_REVERSE_COPY              = C.SEC_ELF_REVERSE_COPY
	SEC_COFF_SHARED                   = C.SEC_COFF_SHARED
	SEC_ELF_COMPRESS                  = C.SEC_ELF_COMPRESS
	SEC_TIC54X_BLOCK                  = C.SEC_TIC54X_BLOCK
	SEC_ELF_RENAME                    = C.SEC_ELF_RENAME
	SEC_TIC54X_CLINK                  = C.SEC_TIC54X_CLINK
	SEC_MEP_VLIW                      = C.SEC_MEP_VLIW
	SEC_COFF_NOREAD                   = C.SEC_COFF_NOREAD
)

const (
	NO_FLAGS             = C.BFD_NO_FLAGS
	HAS_RELOC            = C.HAS_RELOC
	EXEC_P               = C.EXEC_P
	HAS_LINENO           = C.HAS_LINENO
	HAS_DEBUG            = C.HAS_DEBUG
	HAS_SYMS             = C.HAS_SYMS
	HAS_LOCALS           = C.HAS_LOCALS
	DYNAMIC              = C.DYNAMIC
	WP_TEXT              = C.WP_TEXT
	D_PAGED              = C.D_PAGED
	IS_RELAXABLE         = C.BFD_IS_RELAXABLE
	TRADITIONAL_FORMAT   = C.BFD_TRADITIONAL_FORMAT
	IN_MEMORY            = C.BFD_IN_MEMORY
	LINKER_CREATED       = C.BFD_LINKER_CREATED
	DETERMINISTIC_OUTPUT = C.BFD_DETERMINISTIC_OUTPUT
	COMPRESS             = C.BFD_COMPRESS
	DECOMPRESS           = C.BFD_DECOMPRESS
	PLUGIN               = C.BFD_PLUGIN
	COMPRESS_GABI        = C.BFD_COMPRESS_GABI
	FLAGS_SAVED          = C.BFD_FLAGS_SAVED
)

type Architecture C.enum_bfd_architecture

const (
	ArchUnknown Architecture = C.bfd_arch_unknown
	ArchObscure Architecture = C.bfd_arch_obscure
	ArchM68k    Architecture = C.bfd_arch_m68k
	ArchVAX     Architecture = C.bfd_arch_vax
	ArchI960    Architecture = C.bfd_arch_i960
	ArchOr1k    Architecture = C.bfd_arch_or1k
	ArchSparc   Architecture = C.bfd_arch_sparc
	ArchSPU     Architecture = C.bfd_arch_spu
	ArchI386    Architecture = C.bfd_arch_i386
	ArchL1OM    Architecture = C.bfd_arch_l1om
	ArchK1OM    Architecture = C.bfd_arch_k1om
	ArchWe32k   Architecture = C.bfd_arch_we32k
	ArchTahoe   Architecture = C.bfd_arch_tahoe
	ArchI860    Architecture = C.bfd_arch_i860
	ArchI370    Architecture = C.bfd_arch_i370
	ArchRomp    Architecture = C.bfd_arch_romp
	ArchConvex  Architecture = C.bfd_arch_convex
	ArchM88K    Architecture = C.bfd_arch_m88k
	ArchM98K    Architecture = C.bfd_arch_m98k
	ArchPyramid Architecture = C.bfd_arch_pyramid
	ArchH8300   Architecture = C.bfd_arch_h8300
)
