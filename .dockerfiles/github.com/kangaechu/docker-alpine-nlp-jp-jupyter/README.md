docker-alpine-nlp-jp-jupyter
====

Alpine Linuxの自然言語処理Python+Jupyter Notebook実行パッケージ Dockerfile

## Description

Dockerで日本語の自然言語処理関連ツールをPython+Jupyter Notebookで実行するためのDockerfileです。

Docker HubからPullできます。
```
$ docker pull kangaechu/docker-alpine-nlp-jp-jupyter
```
https://hub.docker.com/r/kangaechu/nlp-jp-jupyter/

### 対応しているパッケージ

- [MeCab (+ IPADIC)](http://taku910.github.io/mecab/)
- [Juman](http://nlp.ist.i.kyoto-u.ac.jp/EN/index.php?JUMAN)
- [KNP](http://nlp.ist.i.kyoto-u.ac.jp/?KNP)
- [Juman++](https://github.com/ku-nlp/jumanpp)

- scikit-learn
- numpy
- Jupyter Notebook

## Requirement

Docker

## Usage

以下コマンドを実行します。
イメージが4GBくらいあるので時間がかかります。

```bash
$ docker run -it kangaechu/nlp-jp-jupyter /bin/sh
```

## Licence

MIT

## Author

[kangaechu](https://github.com/kangaechu)
