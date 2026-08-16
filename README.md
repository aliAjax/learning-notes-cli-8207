# learning-notes

一个用 Go 实现的本地 Markdown 学习笔记 CLI，用于快速新建、列出、搜索、编辑和删除学习笔记。每条笔记保存为本地目录中的一个 Markdown 文件，元数据以 JSON front matter 形式写在文件头部，方便后续增加导出、同步等功能。

## 目录结构

```text
.
├── cmd
│   └── learning-notes
│       └── main.go              # 命令入口
├── internal
│   ├── command                  # 命令处理与帮助信息
│   │   ├── command.go
│   │   ├── new.go
│   │   ├── list.go
│   │   ├── search.go
│   │   ├── edit.go
│   │   ├── delete.go
│   │   └── help.go
│   ├── config                   # 数据目录配置
│   │   └── config.go
│   ├── model                    # 笔记模型
│   │   └── note.go
│   ├── search                   # 标题/标签搜索
│   │   └── search.go
│   └── storage                  # Markdown 文件存储
│       └── markdown.go
├── Dockerfile
├── go.mod
└── README.md
```

## 数据目录

笔记默认保存在 `~/.learning-notes`。可以通过以下任一方式修改：

```bash
learning-notes --data-dir /path/to/notes list
learning-notes -d /path/to/notes list
LEARNING_NOTES_DATA_DIR=/path/to/notes learning-notes list
NOTES_DATA_DIR=/path/to/notes learning-notes list
```

配置文件路径由 `--data-dir` 指定；每条笔记对应目录下的 `<id>.md` 文件。全局参数需放在子命令之前，例如：

```bash
learning-notes -d /tmp/notes new --title "Note" --content "Hello"
```

## 构建与运行

### Docker 构建

```bash
docker build -t learning-notes:local .
```

### Docker 运行

容器默认使用 `/notes` 作为数据目录，建议挂载一个本地目录持久化笔记：

```bash
docker run --rm -v "$PWD/data:/notes" learning-notes:local list
```

也可以不挂载卷运行，笔记会保留在容器生命周期内：

```bash
docker run --rm learning-notes:local new \
  --title "Go notes" \
  --tags go,cli \
  --content "Goroutines communicate by sharing memory."
```

## 命令说明

### 全局帮助

```bash
learning-notes --help
learning-notes help
learning-notes help <command>
```

### 新建笔记

```bash
learning-notes new --title <title> [--tags <tags>] (--content <text> | --stdin)
```

示例：

```bash
learning-notes new \
  --title "Go goroutines" \
  --tags go,concurrency \
  --content "Channels provide safe communication."

echo "# Learning log" | learning-notes new --title "Today" --stdin
```

### 列出所有笔记

```bash
learning-notes list
learning-notes list --json
```

### 搜索笔记

自由文本会匹配标题或标签，也可以使用精确筛选参数：

```bash
learning-notes search goroutines
learning-notes search --title "Go goroutines"
learning-notes search --tag go
learning-notes search --tag go --json
```

### 编辑笔记

`--id` 为必填项，至少提供一个要修改的字段：

```bash
learning-notes edit --id <id> --title "New title"
learning-notes edit --id <id> --tags go,notes
learning-notes edit --id <id> --content "Updated content"
echo "Updated content" | learning-notes edit --id <id> --stdin
```

### 删除笔记

```bash
learning-notes delete --id <id>
```

## Markdown 文件格式

每条笔记包含 JSON front matter 和 Markdown 正文：

```markdown
---
{
  "id": "20260816-120000-abcd1234",
  "title": "Go goroutines",
  "tags": ["go", "concurrency"],
  "created_at": "2026-08-16T12:00:00Z",
  "updated_at": "2026-08-16T12:00:00Z"
}
---

Channels provide safe communication.
```
