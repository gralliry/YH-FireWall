---
name: fail-fast-style
description: 代码风格偏好——优先短路写法减少嵌套，避免多层 if 嵌套
metadata:
  type: feedback
---

代码风格要求：优先使用 fail-fast / 短路写法，前置条件不满足就 `continue` / `return` 退出，主路径平铺，不要多层 `if` 嵌套。

例如，下面这种三层 if 嵌套要避免：

```go
for _, addr := range addrs {
    if ipnet, ok := addr.(*net.IPNet); ok {
        if ip, ok := netip.AddrFromSlice(ipnet.IP); ok {
            ...
        }
    }
}
```

应改成：

```go
for _, addr := range addrs {
    ipnet, ok := addr.(*net.IPNet)
    if !ok {
        continue
    }
    ip, ok := netip.AddrFromSlice(ipnet.IP)
    if !ok {
        continue
    }
    ...
}
```

**Why:** 嵌套越少，代码路径越清晰。每个前置条件独立一行，失败立即退出，主逻辑不缩在多层分支里。

**How to apply:** 写代码时优先考虑 "如果不满足条件就提前退出"，而不是 "如果满足条件就缩进执行"。
