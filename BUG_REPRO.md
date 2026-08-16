# Bug 复现说明

## Bug 是什么
存储层在笔记不存在时没有返回 `storage.ErrNotFound` sentinel，命令层因此无法识别“未找到”语义，错误类型在包装时也丢失。

## 如何触发
1. 对不存在的数据目录或笔记执行 `delete --id missing` 或 `edit --id missing`。
2. 在测试中用 `errors.Is(err, storage.ErrNotFound)` 检查返回值。

## 错误信息
返回普通错误文本 `note "missing" not found`，但 `errors.Is(err, storage.ErrNotFound)` 为 false。
