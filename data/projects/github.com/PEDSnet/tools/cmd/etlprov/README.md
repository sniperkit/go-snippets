# PEDSnet ETL Provenance Validator

## Download

Download the latest release on the [releases page](https://github.com/PEDSnet/tools/releases).

## Usage

```bash
$ pedsnet-etlprov [-model <model>] [-version <version>] [-truncate=false] [-ignore <entities>] [-service <service>] [-v] <dir>
```

Option | Description
---|---
-model <model> | Specify a data model other than the default of 'pedsnet', e.g. 'i2b2_pedsnet'
-version <version> | Specify a model version other than the default of '2.0.0', e.g. 2.1.0 or 2.2.0
-truncate=false | Show all errors, even redundant or excessive ones
-service <service> | Specify a model service other than http://data-models.origins.link (not useful unless you run your own model service)
-ignore <entities> | Ignore a comma-separated list of entities (not normally used)
-v | Show the version of this tool

### Example

```bash
$ pedsnet-etlprov -version 2.1.0 -truncate=false
Validating against model 'pedsnet/2.1.0'
Scanning files in '.'
---
1 error has been detected for 'steps.csv'
* Step 0 does not exist
---
1 error has been detected for 'tools.csv'
* High step in range not defined 55
---
1 entity does not have steps
* observation_period
---
197 entities
49 steps
0 tools
3 sources
3 persons
```
