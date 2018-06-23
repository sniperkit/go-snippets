# elastic-examples

```
cd src/github.com/olivere
git clone git@github.com:stormasm/elastic-examples.git
cd elastic-examples/bulkstring
go get golang.org/x/sync/errgroup
go install
bulkstring -index=warehouse -type=product -n=100 -bulk-size=10
```
