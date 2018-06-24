# cmd/ingester

## What is an Ingester?

Ingesters are binaries which periodically send new documents from their respective sources to the Omniscience ingestion service.
An ingester can be written by any one in any language as long as it can communuicate with the ingestion service via gRPC.

## Golang cmd/ingester

The stock ingesters which live under the cmd/ingester package of Omniscience all extend from the common command line
interface provided by common/cli.go. This provides a common set of flags, error handling and ingestion service communication code.

The flags accepted by all ingesters in cmd/ingester are

* --ingestion_service_address
* --modified_since

Ingesters may accept their own custom flags, refer to the README.md in their respective directories.

## How I implement a new ingester in cmd/ingester

An implementation of a new ingester consists of 3 parts.

1.  An implementation of the DocumentFetcher interface.
2.  A DocumentFetcherFactory (a function which returns an instance of the DocumentFetcher)
3.  A main method, which passes the DocumentFetcherFactory to the CreateIngesterCLI method in the common package.

Check out the Google Drive ingester for inspiration.
