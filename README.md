# MySQL to Markdown 文档生成工具

[![Go Report Card](https://goreportcard.com/badge/github.com/jayecc/mysql2md)](https://goreportcard.com/report/github.com/jayecc/mysql2md)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

mysql2md 是一个简单易用的命令行工具，可以将 MySQL 数据库表结构导出为 Markdown 格式的文档，便于查看和分享数据库设计。

## 功能特性

- 🔄 将 MySQL 数据库表结构转换为 Markdown 格式
- 📄 支持生成单个文件或多个文件
- 🔧 可选择是否包含表的 DDL 语句
- 📁 支持自定义输出目录
- ⚡ 支持并发处理提高效率
- 📊 显示表的基本信息、字段详情

## 安装

### 使用 Go 安装（推荐）

```bash
go install github.com/jayecc/mysql2md@latest
```

### 从源码构建

```bash
git clone https://github.com/jayecc/mysql2md.git
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

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。