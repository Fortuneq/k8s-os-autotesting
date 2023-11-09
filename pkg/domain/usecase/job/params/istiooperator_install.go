package params

type IstioOperatorInstallParams struct {
	// Путь до файла IstioOperator
	Path string
	// Ресурс IsitoOperator в виде строки. Имеет больший приоритет, чем Path
	Resource string
	// Время ожидания успешности установки IstioOperator. IstioOperator считается успешно установленным, если Status имеет значеие HEALTHY
	Timeout string
	// Парамерты для подстановки в темплейт
	Values string
}
