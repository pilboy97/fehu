package core

type Config struct {
	Currency string
}

func LoadConfig() error {
	var config Config
	var err error

	if config, err = GetConfig(); err != nil {
		stmt := `insert into config(key, currency) values(?,?)`
		_, err = DB.Exec(stmt, "currency", Code)

		return err
	}

	Code = config.Currency

	return nil
}
func GetConfig() (Config, error) {
	MustDB()

	var ret Config
	stmt := `select currency from config`
	err := DB.QueryRow(stmt).Scan(&ret.Currency)
	if err != nil {
		return Config{}, err
	}

	return ret, nil
}
func SetConfig(config Config) error {
	MustDB()

	stmt := `update config set currency=?`
	_, err := DB.Exec(stmt, config.Currency)
	return err
}
