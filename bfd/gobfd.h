bfd_vma getStartAddress(bfd *abfd);
void *readSymbolTable(bfd *abfd, bfd_boolean dynamic, unsigned int *size, long *count);
int getFileFlags(bfd *abfd);
void setFileFlags(bfd *abfd, int flags);
long getSymtabUpperBound(bfd *abfd);
long getDynamicSymtabUpperBound(bfd *abfd);
long canonicalizeSymtab(bfd *abfd, asymbol **syms);
long canonicalizeDynamicSymtab(bfd *abfd, asymbol **syms);
enum bfd_flavour getFlavor(bfd *abfd);
bfd_boolean findNearestLineDiscriminator(
    bfd *        abfd,
    asection *   sections,
    asymbol **   syms,
    bfd_vma      addr,
    const char **filename, const char **function, unsigned int *line, unsigned int *discriminator);
void printfVMA(bfd *abfd, bfd_vma vma);
bfd_size_type getSectionSize(asection *section);
bfd_vma getSectionVMA(bfd *abfd, asection *section);
void mapOverSections(bfd *abfd, asection *section, void *data);
void goMapOverSections(bfd *abfd, asection *section, void *data);
bfd_boolean findInlinerInfo(bfd *abfd, const char **filename, const char **function, uint *line);