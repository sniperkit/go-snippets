# xmp

<!-- Repo Badges for: Github Project, Slack, License-->

[![GitHub](https://img.shields.io/badge/project-Data_Together-487b57.svg?style=flat-square)](http://github.com/datatogether)
[![Slack](https://img.shields.io/badge/slack-Archivers-b44e88.svg?style=flat-square)](https://archivers-slack.herokuapp.com/)
[![License](https://img.shields.io/github/license/datatogether/xmp.svg)](./LICENSE) 

XMP is a package for parsing Extensible Metadata Platform documents. This
package includes lots of comments to help make sense of XMP for the purposes of
metadata extraction & conversion to other metadata formats.

## License & Copyright

Copyright (C) 2017 Data Together  
This program is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation, version 3.0.

This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE.

See the [`LICENSE`](./LICENSE) file for details.

## Getting Involved

We would love involvement from more people! If you notice any errors or would like 
to submit changes, please see our [Contributing Guidelines](./.github/CONTRIBUTING.md). 

We use GitHub issues for [tracking bugs and feature requests](https://github.com/datatogether/xmp/issues) 
and Pull Requests (PRs) for [submitting changes](https://github.com/datatogether/xmp/pulls)

## Usage

Add to any Golang project with:

`import "github.com/archivers-space/xmp"`

For more information on what XMP is and what it does, read the 
[XMP Wikipedia page](https://en.wikipedia.org/wiki/Extensible_Metadata_Platform): 

> The Extensible Metadata Platform (XMP) is an ISO standard (ISO 16684-1), originally created by Adobe Systems Inc., for the creation, processing and interchange of standardized and custom metadata for digital documents and data sets. XMP standardizes a data model, a serialization format and core properties for the definition and processing of extensible metadata. It also provides guidelines for embedding XMP information into popular image, video and document file formats, such as JPEG and PDF, without breaking their readability by applications that do not support XMP. Therefore, the non-XMP metadata have to be reconciled with the XMP properties.

Technical documentation can be built with `godoc .` or, if your `$GOPATH` and repo 
structure is set up correctly, with something like `godoc -http=:6060 &` followed 
by browsing to http://localhost:6060/pkg/github.com/datatogether.

## Development

More information coming soon; in the meantime please feel free to file issues and 
to improve this README via pull requests.
