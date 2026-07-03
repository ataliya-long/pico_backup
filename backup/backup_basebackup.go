package backup


type basebackup_options struct{
	Host string  // 主机地址
	Port int     // 端口
	User string  // 用户名
	Password string // 密码
	Database string // 数据库名称
	BackupDir string // 备份目录
}

type basebackup interface {
	// 执行备份操作
	Backup(opts basebackup_options) (basebackpath string, err error)
}

