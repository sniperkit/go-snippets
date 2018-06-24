# Package sticker provides a framework for multi-label classification.
![sticker logo](https://rawgit.com/hiro4bbh/sticker/master/logo.svg)

[![Build Status](https://travis-ci.org/hiro4bbh/sticker.svg?branch=master)](https://travis-ci.org/hiro4bbh/sticker)
[![Report Status](https://goreportcard.com/badge/github.com/hiro4bbh/sticker)](https://goreportcard.com/report/github.com/hiro4bbh/sticker)

Copyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).

# Introduction
Package sticker provides a framework for multi-label classification.

sticker is written in golang, so everyone can easily modify and compile it on almost every environments.
You can see sticker's document on [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker).

# Installation
First, download [golang](https://golang.org/), and install it.
Next, get and install sticker as follows:

```
go get github.com/hiro4bbh/sticker
go install github.com/hiro4bbh/sticker/sticker-util
```

Everything has been installed, then you can try sticker's utility command-line tool `sticker-util` now!

# Prepare Datasets
First of all, you should prepare datasets.
sticker assumes the following directory structure for a dataset:

```
+ dataset-root
|-- train.txt: training dataset
|-- text.txt: test dataset
|-- feature_map.txt: feature map (optional)
|-- label_map.txt: label map (optional)
```

Training and test datasets must be formatted as `ReadTextDataset` can handle (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker#ReadTextDataset) for data format).
Feature and label maps should enumerate the name of each feature and label per line in order of identifier, respectively.

You can check the summary of the dataset at `localhost:8080/summary` as follows (you can change the port number with option `addr`):

```
sticker-util -verbose -debug <dataset-root> @summarize -table=<table-filename-relative-to-root>
```

If `featureMap` and `labelMap` is empty string, then feature and label maps are ignored, respectively.

# Implemented Models
## `LabelNearest`: Sparse Weighted Nearest-Neighbor Method
`LabelNearest` is _Sparse Weighted Nearest-Neighbor Method_ __(Aoshima+ 2018)__ which achieved SOTA performances on several XMLC datasets __(Bhatia+ 2016)__.
_Recently, the model can process each data entry faster in 15.1 (AmazonCat-13K), 1.14 (Wiki10-31K), 4.88 (Delicious-200K), 15.1 (WikiLSHTC-325K), 4.19 (Amazon-670K), and 15.5 ms (Amazon-3M) on average, under the same settings of the paper (compare to the original result)._

For example, you can test this method on Amazon-3M dataset __(Bhatia+ 2016)__ as follows:

```
sticker-util -verbose -debug ./data/Amazon-3M/ @trainNearest @testNearest -S=75 -alpha=2.0 -beta=1
```

See the help of `@trainNearest` and `@testNearest` for the sub-command options.

## `LabelNear`: A faster implementation of `LabelNearest`
`LabelNear` is a faster implementation of `LabelNearest` which uses the optimal Densified One Permutation Hashing (DOPH) and the reservoir sampling.
This method can process every data entry in about 1 ms with little performance degradation.
You can see the results on several XMLC datasets __(Bhatia+ 2016)__ at [Dropbox](https://www.dropbox.com/sh/zjerizvew765t0p/AACra7LB0EFwK3RNbSZNprUia?dl=0).

Almost parameters and options are same with the ones of `LabelNearest`.
See the help of `@trainNear` and `@testNear` for details.

## Other Models
### Implemented in core
- `LabelConst`: Multi-label constant model (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker#LabelConst))
- `LabelOne`: One-versus-rest classifier for multi-label ranking (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker#LabelOne))

### Implemented in plugin
- `LabelBoost`: Multi-label Boosting model (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker/plugin#LabelBoost))
- `LabelForest`: Variously-modified FastXML model (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker/plugin#LabelForest))
- `LabelNext`: Your next-generation model (you can add your own train and test commands, see [plugin/next/init.go](https://github.com/hiro4bbh/sticker/blob/master/plugin/next/init.go))

# Implemented Binary Classifiers
## In core (recommended)
- `L1Logistic_PrimalSGD`: L1-logistic regression with stochastic gradient descent (SGD) solving the primal problem (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker#BinaryClassifierTrainer_L1Logistic_PrimalSGD))
- `L1SVC_PrimalSGD`: L1-Support Vector Classifier with SGD solving the primal problem (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker#BinaryClassifierTrainer_L1SVC_PrimalSGD))

## In plugin (not-recommended; for comparison only)
- `L1SVC_DualCD`: L1-Support Vector Classifier with coordinate descent (CD) solving the dual problem (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker/plugin#BinaryClassifierTrainer_L1SVC_DualCD))
- `L2SVC_PrimalCD`: L2-Support Vector Classifier with CD solving the primal problem (see [GoDoc](https://godoc.org/github.com/hiro4bbh/sticker/plugin#BinaryClassifierTrainer_L2SVC_PrimalCD))

# References
- __(Aoshima+ 2018)__ T. Aoshima, K. Kobayashi, and M. Minami. "Revisiting the Vector Space Model: Sparse Weighted Nearest-Neighbor Method for Extreme Multi-Label Classification." [arXiv:1802.03938](https://arxiv.org/abs/1802.03938), 2018.
- __(Bhatia+ 2016)__ K. Bhatia, H. Jain, Y. Prabhu, and M. Varma. The Extreme Classification Repository. 2016. Retrieved January 4, 2018 from [http://manikvarma.org/downloads/XC/XMLRepository.html](http://manikvarma.org/downloads/XC/XMLRepository.html)
