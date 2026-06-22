# pdf_ng

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
pdf_ng input.pdf
```

デフォルトでは`input/`のようにPDF名のフォルダを作り、複数ページのPNGをその中にまとめて出力します。変換中は次のようなプログレスバーが表示されます。

```text
[###############---------------]  50% (2/4)
```

出力先やDPIを指定する場合:

```sh
pdf_ng -o out --dpi 200 input.pdf
```

並列数を指定する場合:

```sh
pdf_ng --jobs 4 input.pdf
```

デフォルトではCPU数を使ってページ単位で並列変換します。

ページ範囲を指定する場合:

```sh
pdf_ng -o out --first 2 --last 5 input.pdf
```

出力ファイル名のprefixを指定する場合:

```sh
pdf_ng -o out --prefix page input.pdf
```

`pdftoppm`の仕様に合わせて、`out/page-1.png`のようなファイルが生成されます。

複数PDFを一括変換する場合:

```sh
pdf_ng *.pdf
```

`-o out`を併用した場合は、`out/input/input-1.png`のようにPDFごとのサブディレクトリへ出力します。

途中で失敗したPDFがあっても残りを変換する場合:

```sh
pdf_ng --continue-on-error *.pdf
```

サブコマンド形式でも実行できます。

```sh
pdf_ng convert -o out input.pdf
```

## インストール

ローカルでインストールする場合:

```sh
go install ./cmd/pdf_ng
```

`$(go env GOPATH)/bin`にPATHが通っていれば、どこからでも`pdf_ng`として実行できます。

GitHubからインストールする場合:

```sh
go install github.com/HizKz/pdf_ng/cmd/pdf_ng@latest
```

## ビルド

```sh
go build -o pdf_ng ./cmd/pdf_ng
```

複数OS向けのバイナリを作る場合:

```sh
./scripts/build-release.sh
```
