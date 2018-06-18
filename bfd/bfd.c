#include <bfd.h>
#include <stdio.h>
#include <stdlib.h>
#include "gobfd.h"

bfd_vma
getStartAddress(bfd *abfd)
{
	return bfd_get_start_address(abfd);
}

void *
readSymbolTable(bfd *abfd, bfd_boolean dynamic, unsigned int *size, long *count)
{
	void *syms;
	*count = bfd_read_minisymbols(abfd, dynamic, &syms, size);
	return syms;
}

int
getFileFlags(bfd *abfd)
{
	return bfd_get_file_flags(abfd);
}

void
setFileFlags(bfd *abfd, int flags)
{
	bfd_set_file_flags(abfd, flags);
}

long
getSymtabUpperBound(bfd *abfd)
{
	return bfd_get_symtab_upper_bound(abfd);
}

long
getDynamicSymtabUpperBound(bfd *abfd)
{
	return bfd_get_dynamic_symtab_upper_bound(abfd);
}

long
canonicalizeSymtab(bfd *abfd, asymbol **syms)
{
	return bfd_canonicalize_symtab(abfd, syms);
}

long
canonicalizeDynamicSymtab(bfd *abfd, asymbol **syms)
{
	return bfd_canonicalize_dynamic_symtab(abfd, syms);
}

enum bfd_flavour
getFlavor(bfd *abfd)
{
	return bfd_get_flavour(abfd);
}

bfd_boolean
findNearestLineDiscriminator(
    bfd *        abfd,
    asection *   sections,
    asymbol **   syms,
    bfd_vma      addr,
    const char **filename, const char **function, unsigned int *line, unsigned int *discriminator)
{
	return bfd_find_nearest_line_discriminator(abfd, sections, syms, addr, filename, function, line, discriminator);
}

void
printfVMA(bfd *abfd, bfd_vma vma)
{
	bfd_printf_vma(abfd, vma);
}

bfd_size_type
getSectionSize(asection *section)
{
	return bfd_get_section_size(section);
}

bfd_vma
getSectionVMA(bfd *abfd, asection *section)
{
	return bfd_get_section_vma(abfd, section);
}

void
mapOverSections(bfd *abfd, asection *section, void *data)
{
	goMapOverSections(abfd, section, data);
}

bfd_boolean
findInlinerInfo(bfd *abfd, const char **filename, const char **function, uint *line)
{
	return bfd_find_inliner_info(abfd, filename, function, line);
}