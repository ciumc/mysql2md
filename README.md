# MySQL to Markdown 文档生成工具

[![Go Report Card](https://goreportcard.com/badge/github.com/ciumc/mysql2md)](https://goreportcard.com/report/github.com/ciumc/mysql2md)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

mysql2md 是一个简单易用的命令行工具，可以将 MySQL 数据库表结构导出为 Markdown 格式的文档，便于查看和分享数据库设计。

## 功能特性

- 将 MySQL 数据库表结构转换为 Markdown 格式
- 支持生成单个文件或多个文件
- 可选择是否包含表的 DDL 语句
- 支持自定义输出目录
- 并发查询提高效率，顺序写入保证输出一致性
- 自动转义 Markdown 特殊字符，避免表格渲染异常
- 显示表的基本信息、字段详情

## 安装

### 使用 Go 安装

```bash
go install github.com/ciumc/mysql2md@latest
```

### 从源码构建

```bash
git clone https://github.com/ciumc/mysql2md.git
cd mysql2md
go build -o mysql2md
```

## 使用方法

### 基本语法

```bash
mysql2md [选项]
```

### 命令行选项

| 选项 | 类型 | 默认值 | 描述 |
|------|------|--------|------|
| `-dsn` | string | `"username:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s"` | 数据库连接字符串 |
| `-dir` | string | `"./output"` | 文档输出目录 |
| `-whole` | bool | `false` | 是否生成单个文件（true: 单个文件，false: 多个文件） |
| `-ddl` | bool | `false` | 是否包含表的 DDL 语句 |

### DSN 连接字符串格式

```
username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s
```

### 使用示例

1. 生成包含 DDL 语句的单个文档文件到当前目录：
   ```bash
   mysql2md -dsn 'root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s' -whole -ddl -dir=.
   ```

2. 生成包含 DDL 语句的单个文档文件：
   ```bash
   mysql2md -dsn 'root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s' -whole -ddl
   ```

3. 为每个表生成单独的文件并包含 DDL 语句：
   ```bash
   mysql2md -dsn 'root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s' -ddl
   ```

4. 为每个表生成单独的文件但不包含 DDL 语句：
   ```bash
   mysql2md -dsn 'root:password@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s'
   ```

## 输出格式

工具会生成以下信息：

### 表信息清单文件
默认生成一个表清单文件，包含数据库中所有表的基本信息：

- 表名（带链接）
- 存储引擎
- 创建时间
- 字符集校对规则
- 表注释

### 表详细信息文件
为每个表生成详细信息文档，包含：

#### 字段信息
- 字段名
- 默认值
- 是否可为空
- 字符集校对规则
- 字段类型
- 键类型（主键、索引等）
- 额外信息（自增等）
- 字段注释

#### DDL 语句（可选）
当使用 `-ddl` 参数时，还会包含表的完整创建语句。

## 项目架构

```
mysql2md/
├── main.go              # CLI 入口：flag 解析、编排
├── db.go                # 数据库连接 + 查询（SchemaQuerier 接口实现）
├── models.go            # 数据模型 + SchemaQuerier 接口定义
├── render.go            # Markdown 渲染（纯函数，无 I/O）
├── writer.go            # 文件写入编排（并发查询 + 顺序写入）
├── render_test.go       # 渲染测试
├── writer_test.go       # 写入测试
├── testhelper_test.go   # SQLite 内存数据库测试辅助
├── db_test.go           # 数据库层测试
├── main_test.go         # 入口逻辑测试
├── go.mod
└── go.sum
```

### 核心设计

- **SchemaQuerier 接口** — 解耦数据库访问与 Markdown 生成，便于测试和扩展
- **并发查询 + 顺序写入** — 利用 errgroup 并发查询数据库，结果按索引存入 slice 后顺序写入文件，兼顾性能与输出一致性
- **纯函数渲染** — render.go 中的函数仅接受结构体、返回字符串，无副作用

## 开发

### 运行测试

```bash
go test -race ./...
```

### 查看测试覆盖率

```bash
go test -race -cover ./...
```

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。
