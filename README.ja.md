# Dircard

[English](README.md) | [日本語](README.ja.md)

Dircard は，`.dircard` ファイルに書かれたディレクトリノートを表示する軽量な CLI ツールです．

複数プロジェクトを行き来する際，ローカルな文脈をターミナル上で素早く確認できます．

## 特徴

- カレントディレクトリまたは親ディレクトリから最も近い `.dircard` を探索して表示
- 出力サイズ，行範囲，探索深度を指定可能
- スクリプト連携向けの JSON 出力
- `bash`，`zsh`，`pwsh` のシェル統合
- シェルフックの安全なインストール/アンインストール

## インストール

Go でインストール:

```bash
go install github.com/dircard/dircard@latest
```

## クイックスタート

1. プロジェクトディレクトリに `.dircard` ファイルを作成
2. 表示したいメモを記述
3. 手動で表示を確認:

```bash
dircard show
```

4. シェル統合を有効化:

```bash
dircard install bash
# または
dircard install zsh
# または
dircard install pwsh
```

反映するにはシェルの再起動，または rc ファイルの再読み込みを行ってください．

5. ディレクトリ移動で自動表示を確認

```bash
cd your/project
# ディレクトリ移動で .dircard の内容が表示される
```

## コマンド

### ノート表示

```bash
dircard show
dircard show --full
dircard show --path
dircard show --json
dircard show --depth 3 --start 10 --lines 20
```

- `dircard show` はカレントディレクトリから親方向に最も近い `.dircard` を探索して表示します．
- `--full` はファイル全体を表示します．
- `--path` は `.dircard` ファイルのパスを表示します．
- `--json` はスクリプト利用向けに JSON 形式で出力します．
- `--depth`，`--start`，`--lines` で探索深度と表示範囲を制御できます．

### シェルフックをインストール

```bash
dircard install [bash|zsh|pwsh]
dircard install pwsh --force
```

`--force` は dircard のフックブロックのみを更新し，他の設定内容は上書きしません．

### シェルフックをアンインストール

```bash
dircard uninstall [bash|zsh|pwsh]
dircard uninstall --force
```

シェル引数を省略した場合，対応する全シェルからフックを削除します．

## 開発

ローカルビルド:

```bash
go build ./...
```

開発中に直接実行:

```bash
go run . show
```

利用可能なコマンドやフラグを確認:

```bash
go run . --help
go run . show --help
```

## 著者

yhotta240 [https://github.com/yhotta240](https://github.com/yhotta240)

## ライセンス

MIT
