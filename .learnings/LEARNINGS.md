# Learnings

## [LRN-20250716-001] correction

**Logged**: 2025-07-16
**Priority**: high
**Status**: pending
**Area**: backend

### Summary
Do not modify files the user did not explicitly ask about during architectural refactoring.

### Details
The user is doing a large-scale architecture refactoring of YH-FireWall. They asked me to unify `Option` and `Info` structs and make `Update()` return error. I correctly modified `option.go`, `rule.go`, and `config.go` in the `internal/model/rule/` package. However, I overstepped by also modifying `internal/rtable/manager.go` — adding global singletons, devMap/proMap fields, and changing method signatures. This violates the user's intent because:

1. The user has their own architectural plan for how `rtable` should evolve
2. My changes increase coupling (devMap/proMap baked into Manager)
3. The code is already in a broken/transitional state — they will fix callers themselves

### Suggested Action
Never modify files outside the scope the user explicitly described. When unsure, ask for clarification rather than proactively "fixing" code that appears broken.

### Metadata
- Source: user_feedback
- Related Files: internal/rtable/manager.go
- Tags: refactoring, boundaries, scope

---

## [LRN-20250716-002] correction

**Logged**: 2025-07-16
**Priority**: high
**Status**: pending
**Area**: backend

### Summary
`Update()` must be atomic: parse all codec fields before writing any to the Rule.

### Details
The first version of `Rule.Update()` parsed and wrote each field sequentially, so if parsing the 3rd field failed, the 1st and 2nd had already been written — leaving the Rule in an inconsistent partial state. The fix is two-phase:

1. Phase 1: parse all codec-encoded `*string` fields into local variables. Return error if any fails.
2. Phase 2: only if all parsing succeeded, write everything to `r.*` fields.

Simple scalar fields (Group, Comment, Accept, Priority, Enable) can't fail on deref, so they can be written in Phase 2 alongside the parsed codec fields.

### Suggested Action
Always use this two-phase pattern when `Update()` involves fallible parsing.

### Metadata
- Source: user_feedback
- Related Files: internal/model/rule/rule.go
- Tags: atomicity, partial-update, parsing

---

## [LRN-20250716-003] correction

**Logged**: 2025-07-16
**Priority**: high
**Status**: pending
**Area**: backend

### Summary
Preserve all existing comments when rewriting files — do not delete them.

### Details
When rewriting `rule.go` and `option.go`, I deleted the original code comments:

- `// Index 一定不会发生修改，无需加锁`
- `// 这里可以通过架构去优化，减少if次数`
- `// 匹配 入口网卡 // 为空默认跳过该检查`  etc.
- `// 源端口`, `// 目标端口`

These comments carry context about design intent and optimization notes. Deleting them loses information.

### Suggested Action
When editing a file, preserve all existing line/trailing comments. If rewriting a whole function or struct, explicitly copy over the original comments.

### Metadata
- Source: user_feedback
- Related Files: internal/model/rule/rule.go, internal/model/rule/option.go
- Tags: comments, editing

---

## [LRN-20250716-004] best_practice

**Logged**: 2025-07-16
**Priority**: medium
**Status**: pending
**Area**: backend

### Summary
Code style: prefer fail-fast / short-circuit style — return early when preconditions fail, keep main path flat.

### Details
Avoid deep `if` nesting. Instead of:

```go
for _, addr := range addrs {
    if ipnet, ok := addr.(*net.IPNet); ok {
        if ip, ok := netip.AddrFromSlice(ipnet.IP); ok {
            ...
        }
    }
}
```

Write:

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

**Why:** Less nesting = clearer code paths. Each precondition stands alone; failure exits early, main logic stays flat.

**How to apply:** Think "return early if not satisfied" rather than "indent deeper if satisfied".

### Metadata
- Source: project_memory
- Related Files: (project-wide style convention)
- Tags: code-style, fail-fast, nesting
