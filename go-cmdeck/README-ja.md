# Go CmDeck

Goで書かれたRundeck風CLI/TUIジョブ実行管理ツールです。Go CmDeckは包括的な実行履歴を持つジョブの定義、実行、追跡を可能にし、反復的なタスクやコマンドの管理を簡単にします。

## 機能

- **ジョブ実行管理**: ラベル、説明、コマンド、変数でジョブを定義
- **実行履歴**: タイムスタンプ、終了コード、成功/失敗ステータス、詳細出力でジョブ実行を追跡
- **CLIインターフェース**: ジョブ操作用のコマンドラインインターフェース (list, run, add, remove)
- **TUIインターフェース**: ✓/✗アイコンでジョブステータス表示するインタラクティブターミナルユーザーインターフェース
- **変数置換**: `${VAR}` 構文を使用した環境変数展開でジョブを実行
- **詳細ログ**: STDOUT/STDERR分離付き包括的実行レポート
- **設定可能テーマ**: TUIインターフェース用のカスタマイズ可能カラーテーマ
- **レスポンシブレイアウト**: オーバーフロー処理付き適応型ターミナルレイアウト

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/ToshihitoKon/agent-works.git
cd agent-works/go-cmdeck

# バイナリをビルド
go build -o go-cmdeck

# 実行可能にして、オプションでPATHに移動
chmod +x go-cmdeck
# sudo mv go-cmdeck /usr/local/bin/  # オプション: グローバルインストール
```

## クイックスタート

```bash
# サンプルジョブで初期化
./go-cmdeck init

# 実行ステータス付きジョブ一覧表示
./go-cmdeck list

# 特定のジョブを実行
./go-cmdeck run monitoring

# インタラクティブTUIモード開始
./go-cmdeck tui

# 新しいジョブを追加
./go-cmdeck add -name "backup" -label "データベースバックアップ" -description "日次バックアップジョブ"
```

## コマンド

| コマンド | 説明 |
|---------|-------------|
| `init` | サンプルジョブで設定を初期化 |
| `list`, `ls` | 実行ステータス付き全ジョブ一覧表示 |
| `execute`, `exec <name>` | ジョブを実行して実行履歴を記録 |
| `run <name>` | ジョブを実行して実行履歴を記録 |
| `add` | 新しいジョブを追加（インタラクティブ） |
| `remove`, `rm <name>` | ジョブを削除 |
| `tui` | TUIモードを開始 |
| `help` | ヘルプを表示 |

## TUIインターフェース

TUI（ターミナルユーザーインターフェース）はジョブの管理と実行のためのインタラクティブな方法を提供します：

- **ジョブリスト**: ステータスアイコン付き全ジョブを表示（成功は✓、失敗は✗）
- **ナビゲーション**: 矢印キーまたはj/kを使用してナビゲート
- **ジョブ実行**: スペースキーを押して選択されたジョブを実行
- **ジョブ詳細**: 下部パネルで選択されたジョブの詳細情報を表示
- **リアルタイム更新**: 実行ステータスがリアルタイムで更新

### TUI操作

- `↑/↓` または `j/k`: ジョブ間をナビゲート
- `Space`: 選択されたジョブを実行
- `q` または `Ctrl+C`: 終了

## 設定

設定は `~/.config/go-cmdeck/config.json` に保存されます：

```json
{
  "contexts": {
    "monitoring": {
      "name": "monitoring",
      "label": "システム監視",
      "description": "システム監視ツールを有効化",
      "commands": {
        "run": "echo '監視が有効になりました' && ps aux | head -5"
      },
      "variables": {
        "LOG_PATH": "/var/log/monitoring",
        "INTERVAL": "5"
      },
      "last_result": {
        "timestamp": "2025-06-11T22:56:44.500268+09:00",
        "success": true,
        "exit_code": 0,
        "output": "コマンド実行出力..."
      }
    }
  },
  "theme": {
    "title": "205",
    "selected": "199", 
    "border": "168",
    "output_title": "212"
  }
}
```

### ジョブ構造

各ジョブは以下で構成されます：

- **name**: ジョブの一意識別子
- **label**: 人間が読める表示名
- **description**: ジョブが何をするかのオプション説明
- **commands.run**: 実行するコマンド
- **variables**: 変数置換用のキー値ペア
- **last_result**: 実行履歴（自動管理）

### 変数置換

コマンド内で `${VARIABLE_NAME}` を使用して変数を置換：

```json
{
  "commands": {
    "run": "echo '${HOST}:${PORT}に接続中'"
  },
  "variables": {
    "HOST": "localhost",
    "PORT": "8080"
  }
}
```

## 例

### バックアップジョブの作成

```bash
./go-cmdeck add -name "backup" -label "データベースバックアップ" -description "PostgreSQLデータベースをバックアップ"
```

その後、設定を編集してコマンドを追加：

```json
{
  "name": "backup",
  "label": "データベースバックアップ", 
  "description": "PostgreSQLデータベースをバックアップ",
  "commands": {
    "run": "pg_dump -h ${DB_HOST} -U ${DB_USER} ${DB_NAME} > backup_$(date +%Y%m%d).sql"
  },
  "variables": {
    "DB_HOST": "localhost",
    "DB_USER": "postgres", 
    "DB_NAME": "myapp"
  }
}
```

### 監視ジョブ

```json
{
  "name": "monitoring",
  "label": "システムヘルスチェック",
  "description": "システムリソースとサービスをチェック",
  "commands": {
    "run": "echo 'システムステータス:' && uptime && echo 'ディスク使用量:' && df -h / && echo 'メモリ:' && free -h"
  },
  "variables": {
    "ALERT_EMAIL": "admin@company.com"
  }
}
```

## 開発

### 前提条件

- Go 1.21以降
- カラーサポート付きターミナル

### 依存関係

- `github.com/charmbracelet/bubbletea`: TUIフレームワーク
- `github.com/charmbracelet/lipgloss`: TUI用スタイリング

### ソースからビルド

```bash
git clone https://github.com/ToshihitoKon/agent-works.git
cd agent-works/go-cmdeck
go mod tidy
go build -o go-cmdeck
```

### テスト

```bash
go test ./...
```

## 貢献

1. リポジトリをフォーク
2. フィーチャーブランチを作成
3. 変更を行う
4. 該当する場合はテストを追加
5. プルリクエストを提出

開発ガイドラインとアーキテクチャの決定については[CLAUDE.md](./CLAUDE.md)を参照してください。

## ライセンス

このプロジェクトはAgent Worksコレクションの一部です。ライセンス情報については親リポジトリを確認してください。

## 類似プロジェクト

- [Rundeck](https://www.rundeck.com/): エンタープライズジョブスケジューラーおよびランブック自動化
- [Ansible](https://www.ansible.com/): IT自動化プラットフォーム
- [Jenkins](https://www.jenkins.io/): CI/CD自動化サーバー

Go CmDeckは個人および小規模チームの使用例向けの軽量でターミナルベースの代替として設計されています。