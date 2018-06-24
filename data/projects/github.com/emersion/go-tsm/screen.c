#include <libtsm.h>

void get_screen_attr_bitfield(struct tsm_screen_attr *attr, bool *bold,
		bool *underline, bool *inverse, bool *protect, bool *blink) {
	*bold = attr->bold;
	*underline = attr->underline;
	*inverse = attr->inverse;
	*inverse = attr->inverse;
	*protect = attr->protect;
	*blink = attr->blink;
}

void set_screen_attr_bitfield(struct tsm_screen_attr *attr, bool bold,
		bool underline, bool inverse, bool protect, bool blink) {
	attr->bold = bold;
	attr->underline = underline;
	attr->inverse = inverse;
	attr->protect = protect;
	attr->blink = blink;
}
