package config

// var configPath string

// func GetConfigName() string {
// 	return configPath
// }

// func CreateConfig() error {
// 	homeDir, err := os.UserHomeDir()
// 	if err != nil {
// 		return err
// 	}

// 	configPath = filepath.Join(homeDir, ".config")
// 	if _, err := os.Stat(configPath); os.IsNotExist(err) {
// 		if err := os.Mkdir(configPath, 0655); err != nil {
// 			return err
// 		}
// 	}

// 	configFileName := "netools.csv"
// 	configPath = filepath.Join(configPath, configFileName)
// 	if _, err := os.Stat(configPath); os.IsNotExist(err) {
// 		f, err := os.Create(configPath)
// 		if err != nil {
// 			return err
// 		}
// 		f.Close()
// 	}

// 	return nil
// }
