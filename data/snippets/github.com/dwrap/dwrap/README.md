# dwrap

`dwrap` is `Docker Wrapper`

**現在はまだ実験的な実装です。**

## 概要

alpineベースのDockerコンテナを動的に生成して指定のコマンドを実行します。

docker-machineを用いることで、リモートマシン上でコマンドを実行させることも可能です。


## 使い方
 
```bash:dwrapの使い方

# alpineベースのコンテナを生成し、そこでcurlコマンドを実行
$ dwrap curl --help 

# jqコマンドを実行、入力としてjsonファイルをパイプで渡す
$ cat samples/jq.json | dwrap jq "."

#----------------------------
# リモートDockerでの例
#----------------------------
# docker-machineでリモートのDockerマシン上でコマンド実行
# 例：重いファイルをダウンロード、圧縮して手元マシンへ

$ eval $(docker-machne env remote-machine)
$ dwrap curl -Ls http://cloud.centos.org/centos/7/atomic/images/CentOS-Atomic-Host-7-Installer.iso | dwrap zip > ~/test.zip

```

## インストール

[リリースページ](https://github.com/dwrap/cli/releases/latest)から各プラットフォームごとにバイナリをダウンロードして展開してください。

## 内部動作

1. alpineをベースイメージに指定するDockerfileを作成
1. 指定のコマンドをapkコマンドでインストール
1. コンテナのentrypointを指定のコマンドに設定

生成されるDockerイメージは`dwrap-image/xxxxx`という名前で作成されます。
一度作成されたイメージはキャッシュされます。

作成されたイメージの一括削除をする場合は以下のコマンドを実行してください。

```
$ docker rmi -f `docker images -aq dwrap-image/*`
```

# License

This project is published under [Apache 2.0 License](LICENSE).

# Author

* Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
