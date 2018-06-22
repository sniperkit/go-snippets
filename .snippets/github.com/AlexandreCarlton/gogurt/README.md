# gogurt

[![Build Status](https://travis-ci.org/AlexandreCarlton/gogurt.svg?branch=master)](https://travis-ci.org/AlexandreCarlton/gogurt)

Downloads and builds static packages from source.

The successor to [`install-from-source`](https://github.com/AlexandreCarlton/install-from-source),
packaged in one, neat binary.

To install it:

```bash
  go install github.com/alexandrecarlton/gogurt/cmd/gogurt
```

## Caveats
For now, this will ignore `musl`, and install `gcc` to get `libstdc++.a`.
To this end, we use a Docker container to build these applications.
