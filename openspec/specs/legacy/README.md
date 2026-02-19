# Legacy Specs (归档规格)

本目录包含旧格式的规格文件，已归档保留作为历史记录。

## 格式说明

这些文件使用**编号文件**格式（如 `001-campaign-management.md`），已被**capability 目录格式**取代。

## 当前格式

OpenSpec 现在使用 capability 目录格式：

```
openspec/specs/
├── campaign-management/
│   └── spec.md
├── feedback-system/
│   └── spec.md
└── ...
```

## 迁移状态

| 旧文件 | 新位置 | 状态 |
|--------|--------|------|
| 001-campaign-management.md | campaign-management/ | ✅ 已迁移 |
| 002-order-payment-system.md | order-payment-system/ | ✅ 已迁移 |
| 003-reward-system.md | - | ⚠️ 待创建 capability |
| 004-sync-adapter.md | - | ⚠️ 待创建 capability |
| 005-mobile-landing-page.md | - | ⚠️ 待创建 capability |

## 注意事项

- 这些文件仅作为历史记录保留
- CLI 工具不识别此格式
- 新的规格请使用 capability 目录格式
- 参考 `openspec/project.md` 和 `openspec/AGENTS.md` 了解当前规范

## 归档时间

2026-02-19
