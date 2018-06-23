
function get_if_stats() {
  printf "%s", $(ifstat -bi en0 0.1 1 | tail -n 1 | \
    awk '
      {printf("↑ %d%s",
        ($1 > 1000) ? $1 / 1024 : $1,
        ($1 > 1000) ? "m" : "k",
        ($1 > 1000) ? "A13D63" : \
          ($1 > 4000) ? "2A9D8F" : \
            ($1 > 2000) ? "A5C882" : "324376")
      }');
}
