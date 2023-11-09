package params

// Параметры для Image Pull Secret
type PullSecretParams struct {
	// Имя секрета
	Name string `yaml:"name"`
	// Неймспейс секрета
	Namespace string `yaml:"namespace"`
	// Имя пользователя
	User string `yaml:"user"`
	// Пароль
	Password string `yaml:"password"`
	// Адрес image registry сервера
	Server string `yaml:"server"`
	// Список ServiceAccounts, для которых нужно добавить созданный секрет в imagePullSecrets
	ServiceAccounts []string `yaml:"serviceAccounts"`
}
