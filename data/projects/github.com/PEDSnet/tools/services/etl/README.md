# ETL Conventions Service

Service that exposes endpoints for fetching and parsing ETL conventions documents used by PEDSnet. The current documents are located here:

- [PEDSnet](https://github.com/PEDSnet/Data_Models/blob/master/PEDSnet/docs/Pedsnet_CDM_ETL_Conventions.md)
- [i2b2](https://github.com/PEDSnet/Data_Models/blob/master/i2b2/V2/docs/i2b2_pedsnet_v2_etl_conventions.md)

## Run

Options:

- `host`
- `port`
- `token`, alternately the `GITHUB_AUTH_TOKEN` environment variable can be set.

```
docker build -t pedsnet/etlconv .

# Token required since the file is in a private repository.
docker run pedsnet/etlconv --token=GITHUB_TOKEN
```

## Usage

Perform a GET request to the root endpoint. See `tests/expected.json` for example output.


## API

### Models

- `pedsnet`
- `i2b2`

### GET /<model>

Returns the parsed document for the specified model.

#### Parameters

- `ref` - Get the document at a specific commit.
- `asof` - Get the document as of some time. Supported formats are:
    - `2015-09-04T10:40:03Z` - ISO 8601 UTC time (with the Z)
    - `2015-09-04T10:40:03` - ISO 8601 UTC time (without the Z)
    - `2015-09-04` - ISO 8601 date without a time

### GET /<model>/commits

Returns a set of commits for the document. The document may be located at different paths over time, so the commit information includes the `file_path`.

### GET /<model>/prov

*Experimental*. Returns a set of provenance events.

#### Parameters

- `ref` - Get the document at a specific commit.
- `asof` - Get the document as of some time. Supported formats are:
    - `2015-09-04T10:40:03Z` - ISO 8601 UTC time.
    - `2015-09-04T10:40:03` - ISO 8601 UTC time (without the Z)
    - `2015-09-04` - ISO 8601 date without a time

### GET /<model>/log

*Experimental*. Returns the full log of provenance events for the document from it's origin. This is ordered relative to the document, not by time.
