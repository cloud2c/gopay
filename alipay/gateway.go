package alipay

// 支付宝网关地址（UTF-8 字符集版本）。调用方通过 IsProd 选择，不要手写地址。
const (
	// ProductionGatewayURL 为正式环境网关。
	ProductionGatewayURL = baseUrlUtf8
	// SandboxGatewayURL 为沙箱环境网关（需使用沙箱应用的 AppID 与密钥）。
	SandboxGatewayURL = sandboxBaseUrlUtf8
)

// GatewayURL 按环境返回网关地址：isProd 为 true 返回正式环境，否则返回沙箱。
func GatewayURL(isProd bool) string {
	if isProd {
		return ProductionGatewayURL
	}
	return SandboxGatewayURL
}
