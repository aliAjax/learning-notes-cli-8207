# Bug 复现说明

## Bug 是什么
空标签被处理成 nil slice，导致 JSON 输出为 `null`；同时搜索逻辑遇到 nil 标签会提前返回 false，无标签笔记按标题也搜不到。

## 如何触发
1. 创建不带标签的笔记，运行 `list --json` 或 `search --json`。
2. 对无标签笔记运行 `search <标题关键词>`。

## 错误信息
JSON 中出现 `"tags":null`；按标题搜索无标签笔记时结果为空。
