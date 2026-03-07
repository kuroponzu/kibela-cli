# kibela-cli

AIフレンドリーなKibela操作用CLIツール

## インストール

```bash
go install github.com/kuroponzu/kibela-cli/cmd/kibela@latest
```

またはソースからビルド:

```bash
git clone https://github.com/kuroponzu/kibela-cli.git
cd kibela-cli
make build
```

## 設定

環境変数を設定してください:

```bash
export KIBELA_TOKEN="your-api-token"
export KIBELA_TEAM="your-team-name"
```

## 使い方

### 記事取得

```bash
# IDで取得
kibela get --id <note-id>

# パスで取得
kibela get --path "/notes/12345"

# JSON出力
kibela get --id <note-id> --json
```

### 記事作成

```bash
# インライン指定
kibela create --title "タイトル" --content "本文" --group-id <group-id>

# ファイルから
kibela create --title "タイトル" --content-file ./article.md --group-id <group-id>

# 標準入力から
cat article.md | kibela create --title "タイトル" --group-id <group-id> --stdin

# 下書きとして作成
kibela create --title "タイトル" --content "本文" --group-id <group-id> --draft

# 共同編集を有効に
kibela create --title "タイトル" --content "本文" --group-id <group-id> --co-editing
```

### 記事更新

```bash
# タイトル更新
kibela update --id <note-id> --title "新タイトル"

# 本文更新
kibela update --id <note-id> --content "新しい本文"

# ファイルから更新
kibela update --id <note-id> --content-file ./updated.md

# 標準入力から更新
cat updated.md | kibela update --id <note-id> --stdin
```

### グローバルフラグ

- `--json`: JSON形式で出力
- `--verbose, -v`: 詳細ログを出力
- `--help, -h`: ヘルプを表示

## 開発

```bash
# ビルド
make build

# テスト
make test

# 全プラットフォーム向けビルド
make build-all

# クリーン
make clean
```

## 終了コード

| コード | 説明 |
|--------|------|
| 0 | 成功 |
| 1 | 設定エラー |
| 2 | 認証エラー |
| 3 | 権限エラー |
| 4 | 未発見 |
| 10 | 入力エラー |
| 20 | GraphQLエラー |

## ライセンス

MIT License
