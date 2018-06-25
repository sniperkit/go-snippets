# Vocabulary Validator

Command line interface for validating various aspects of a vocabulary data set.

To get the full set of options for any command, use the `help` command.

```
pedsnet-vocab help <command>
```

## Examples

### Summary

```
$ pedsnet-vocab summary 2014-10-20-pedsnet-concepts.csv
Detected ',' delimiter
Counts
---
1068072	Domains
48	    Vocabularies
1068072	Concepts
```

### Compare

*Omit the `--quiet` flag to have all the changes written to stdout in pairs.*

```
$ pedsnet-vocab compare --quiet 2015-06-12-pedsnet-concepts.csv 2015-11-02-pedsnet-concepts.csv
Summary:
* 168561 Added
* 3 Removed
* 55632 Changed
```
