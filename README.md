# pf2pg

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
pf2pg input.pdf
```

デフォルトでは`input/`のようにPDF名のフォルダを作り、複数ページのPNGをその中にまとめて出力します。変換中は次のようなプログレスバーが表示されます。

```text
[###############---------------]  50% (2/4)
```

出力先やDPIを指定する場合:

```sh
pf2pg -o out --dpi 200 input.pdf
```

並列数を指定する場合:

```sh
pf2pg --jobs 4 input.pdf
```

デフォルトではCPU数を使ってページ単位で並列変換します。

ページ範囲を指定する場合:

```sh
pf2pg -o out --first 2 --last 5 input.pdf
```

出力ファイル名のprefixを指定する場合:

```sh
pf2pg -o out --prefix page input.pdf
```

`pdftoppm`の仕様に合わせて、`out/page-1.png`のようなファイルが生成されます。

複数PDFを一括変換する場合:

```sh
pf2pg *.pdf
```

`-o out`を併用した場合は、`out/input/input-1.png`のようにPDFごとのサブディレクトリへ出力します。

途中で失敗したPDFがあっても残りを変換する場合:

```sh
pf2pg --continue-on-error *.pdf
```

サブコマンド形式でも実行できます。

```sh
pf2pg convert -o out input.pdf
```

## インストール

ローカルでインストールする場合:

```sh
go install ./cmd/pf2pg
```

`$(go env GOPATH)/bin`にPATHが通っていれば、どこからでも`pf2pg`として実行できます。

GitHubからインストールする場合:

```sh
go install github.com/HizumeKazushi/pf2pg/cmd/pf2pg@latest
```

## ビルド

```sh
go build -o pf2pg ./cmd/pf2pg
```

複数OS向けのバイナリを作る場合:

```sh
./scripts/build-release.sh
```
