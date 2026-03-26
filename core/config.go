package core

type Config struct {
	Currency string
}

func InitConfig() {
	var config Config
	var err error

	if config, err = GetConfig(); err != nil {
		stmt := `insert into config(key, currency) values(?,?)`
		_, err = DB.Exec(stmt, "currency", Code)
		if err != nil {
			panic(err)
		}
	}

	if config.Currency != Code {
		Code = config.Currency
	}
}
func GetConfig() (Config, error) {
	ChkDB()

	var ret Config
	stmt := `select currency from config`
	err := DB.QueryRow(stmt).Scan(&ret.Currency)
	if err != nil {
		return Config{}, err
	}

	return ret, nil
}
func SetConfig(config Config) error {
	ChkDB()

	stmt := `update config set currency=?`
	_, err := DB.Exec(stmt, config.Currency)
	return err
}
