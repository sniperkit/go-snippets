

const archiveServices = [
  {
    id:          "5d5583g5-38a8-26d3-be70-c324bf686a87",
    name:        "Internet Archive",
    description: "the internet archive",
    homeUrl:     "https://archive.org",
  },
  {
    id:          "3f5b22g5-37b4-5dc3-be91-c324bf686a87",
    name:        "EOT Nomination Tool",
    description: "",
    homeUrl:     "https://github.com/edgi-govdata-archiving/eot-nomination-tool",
  },
  {
    id:          "4c0122g5-38a8-40b3-be91-c324bf686a87",
    name:        "archivers.space",
    description: "",
    homeUrl:     "https://www.archivers.space",
  },
  {
    id:          "8d7e22g5-38a8-40b3-be91-c324bf686a87",
    name:        "archivers 2.0",
    description: "",
    homeUrl:     "https://alpha.archivers.space",
  }
]

export function archiveService(id) {
  return archiveServices.find((s) => s.id == id);
}