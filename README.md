![alt text](image.png)
# picobackup

picobackup是PostgreSQL 的自动化备份工具，调用 `pg_dump` 导出后自动 gzip 压缩。Go 单二进制，零外部运行时依赖。

## 目录结构

```
db-backup/
├── main.go
├── go.mod
├── config.yaml          → 配置文件（默认值）
├── config/
│   └── config.go
├── backup/
│   ├── backup.go        → 接口定义 + gzip 压缩
│   └── postgres.go      → pg_dump 驱动
└── README.md
```

## 编译

```bash
cd db-backup
go build -o backup.exe .
```

依赖只有 `gopkg.in/yaml.v3`，`go mod tidy` 自动拉取。

## 快速开始

```bash
# PostgreSQL 全量备份
./backup --db postgres --host 192.168.23.129 --user postgres --database postgres --passwd xxxxx
```

输出文件格式：`{数据库名}_{主机}_{时间戳}.sql.gz`

示例：`postgres_192.168.23.129_20260702_223000.sql.gz`

## 命令行参数

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--db` | 否 | `postgres` | 数据库类型，固定为 `postgres` |
| `--host` | 否 | config 中的值 | 数据库主机地址 |
| `--port` | 否 | 5432 | 数据库端口 |
| `--user` | 否 | config 中的值 | 连接用户 |
| `--passwd` | 否 | config 中的值 | 密码 |
| `--database` | 是 | - | 数据库名 |
| `--dir` | 否 | config 中的 `backup_dir` | 备份输出目录 |
| `--config` | 否 | `config.yaml` | 配置文件路径 |

每个 `--` 参数的优先级：命令行 > 配置文件 > 内置默认值。

## 配置文件 config.yaml

```yaml
backup_dir: ./backups

postgres:
  host: 192.168.23.129
  port: 5432
  user: postgres
  password: ""
```

配置好之后，命令行只需传 `--database` 即可：

```bash
./backup --database postgres --passwd xxxxx
```

## 定时执行

### Linux cron

每天凌晨 2:00 执行：

```
0 2 * * * /opt/db-backup/backup --database mydb --passwd xxxxx >> /var/log/backup.log 2>&1
```

每周日凌晨 3:00 全量：

```
0 3 * * 0 /opt/db-backup/backup --database mydb --passwd xxxxx --dir /data/backups
```

### Windows 任务计划

```powershell
$Action = New-ScheduledTaskAction -Execute "D:\db-backup\backup.exe" `
  -Argument "--database mydb --passwd xxxxx --dir D:\backups"
$Trigger = New-ScheduledTaskTrigger -Daily -At 02:00
Register-ScheduledTask -TaskName "DB-Backup-Postgres" -Action $Action -Trigger $Trigger
```

## 前置要求

本机需安装 `pg_dump`（PostgreSQL 客户端工具）。通过 `PGPASSWORD` 环境变量传递密码，避免交互式提示。

```bash
# Ubuntu / Debian
apt install postgresql-client

# CentOS / RHEL
yum install postgresql

# Windows
# 安装 PostgreSQL 时勾选 Command Line Tools，或将 pg_dump.exe 所在目录加入 PATH
```
