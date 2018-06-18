package bytes

// BaseUnit is a type that defines a unit of measure as a distinct dimensions
// from which all other units can be derived.
type BaseUnit = float64

// International Electrotechnical Commission (IEC) base 2 units [IEEE 1541]
// - http://physics.nist.gov/cuu/Units/binary.html
// - https://en.wikipedia.org/wiki/Binary_prefix
const (
	// Byte is the BaseUnit for all binary prefix sizes
	Byte BaseUnit = 1 << (iota * 10)
	// Kibibyte is a multiple of the BaseUnit Byte.
	Kibibyte
	// Mebibyte is a multiple of the BaseUnit Byte.
	Mebibyte
	// Gibibyte is a multiple of the BaseUnit Byte.
	Gibibyte
	// Tebibyte is a multiple of the BaseUnit Byte.
	Tebibyte
	// Pebibyte is a multiple of the BaseUnit Byte.
	Pebibyte
	// Exbibyte is a multiple of the BaseUnit Byte.
	Exbibyte
	// Zebibyte is a multiple of the BaseUnit Byte.
	Zebibyte
	// Yobibyte is a multiple of the BaseUnit Byte.
	Yobibyte
)

const (
	// KiB is the shorthand for Kibibyte
	KiB = Kibibyte
	// MiB is the shorthand for Mebibyte
	MiB = Mebibyte
	// GiB is the shorthand for Gibibyte
	GiB = Gibibyte
	// TiB is the shorthand for Tebibyte
	TiB = Tebibyte
	// PiB is the shorthand for Pebibyte
	PiB = Pebibyte
	// EiB is the shorthand for Exbibyte
	EiB = Exbibyte
	// ZiB is the shorthand for Zebibyte
	ZiB = Zebibyte
	// YiB is the shorthand for Yobibyte
	YiB = Yobibyte
)

// International Systems of Units (SI) base 10 units
// - http://physics.nist.gov/cuu/Units/prefixes.html
const (
	// Kilobyte is a multiple of the BaseUnit Byte.
	Kilobyte BaseUnit = 1000 * Byte
	// Megabyte is a multiple of the BaseUnit Byte.
	Megabyte = 1000 * Kilobyte
	// Gigabyte is a multiple of the BaseUnit Byte.
	Gigabyte = 1000 * Megabyte
	// Terabyte is a multiple of the BaseUnit Byte.
	Terabyte = 1000 * Gigabyte
	// Petabyte is a multiple of the BaseUnit Byte.
	Petabyte = 1000 * Terabyte
	// Zettabyte is a multiple of the BaseUnit Byte.
	Exabyte = 1000 * Petabyte
	// Zettabyte is a multiple of the BaseUnit Byte.
	Zettabyte = 1000 * Exabyte
	// Yottabyte is a multiple of the BaseUnit Byte.
	Yottabyte = 1000 * Zettabyte
)

const (
	// KB is the shorthand for Kilobyte
	KB = Kilobyte
	// MB is the shorthand for Megabyte
	MB = Megabyte
	// GB is the shorthand for Gigabyte
	GB = Gigabyte
	// TB is the shorthand for Terabyte
	TB = Terabyte
	// PB is the shorthand for Petabyte
	PB = Petabyte
	// EB is the shorthand for Exabyte
	EB = Exabyte
	// ZB is the shorthand for Zettabyte
	ZB = Zettabyte
	// YB is the shorthand for Yottabyte
	YB = Yottabyte
)

var prefixes = map[string]float64{
	"KB":  KB,
	"MB":  MB,
	"GB":  GB,
	"TB":  TB,
	"PB":  PB,
	"EB":  EB,
	"ZB":  ZB,
	"YB":  YB,
	"KiB": KiB,
	"MiB": MiB,
	"GiB": GiB,
	"TiB": TiB,
	"PiB": PiB,
	"EiB": EiB,
	"ZiB": ZiB,
	"YiB": YiB,
}
