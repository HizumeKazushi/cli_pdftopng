# p2p

PDFをPNGに変換するGo製CLIです。PDFのレンダリングにはPopplerの`pdftoppm`を使います。

## 必要なもの

macOS:

```sh
brew install poppler
```

Ubuntu/Debian:

```sh
sudo apt-get install poppler-utils
```

## 使い方

```sh
p2p input.pdf
```

デフォルトでは`input/`のようにPDF名のフォルダを作り、複数ページのPNGをその中にまとめて出力します。変換中は次のようなプログレスバーが表示されます。

```text
[###############---------------]  50% (2/4)
```

出力先やDPIを指定する場合:

```sh
p2p -o out --dpi 200 input.pdf
```

ページ範囲を指定する場合:

```sh
p2p -o out --first 2 --last 5 input.pdf
```

出力ファイル名のprefixを指定する場合:

```sh
p2p -o out --prefix page input.pdf
```

`pdftoppm`の仕様に合わせて、`out/page-1.png`のようなファイルが生成されます。

サブコマンド形式でも実行できます。

```sh
p2p convert -o out input.pdf
```

## インストール

ローカルでインストールする場合:

```sh
go install ./cmd/p2p
```

`$(go env GOPATH)/bin`にPATHが通っていれば、どこからでも`p2p`として実行できます。

GitHubからインストールする場合:

```sh
go install github.com/HizumeKazushi/cli_pdftopng/cmd/p2p@latest
```

## ビルド

```sh
go build -o p2p ./cmd/p2p
```

複数OS向けのバイナリを作る場合:

```sh
./scripts/build-release.sh
```
