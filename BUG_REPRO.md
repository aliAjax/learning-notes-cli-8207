# Bug 复现说明

## Bug 是什么
存储层的 Save/Get/List/Delete 方法忽略了 context 取消信号，即使 ctx 已 cancel 也会继续执行。

## 如何触发
1. 创建一个已经 cancel 的 context。
2. 直接调用存储层 Save/Get/List/Delete。

## 错误信息
调用返回 nil 错误或正常结果，而不是 `context.Canceled`。
