package conn

// 定义一个全局变量，其他包都能访问
var GlobalConnManager *ConnManager

// 初始化这个变量，创建一个新的连接管理器实例
func InitGlobalConnManager() {
	GlobalConnManager = NewConnManager()
}
