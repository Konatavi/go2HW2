package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Konatavi/go2HW2/internal/app/apiserver"
	"github.com/joho/godotenv"
)

var (
	configPath        string
	typeConfigFile    string
	configPathDefault string = "configs/apiserver.toml"
)

/*
Добавить в код необходимые блоки, для того, чтобы можно было запускать приложение следующими командами:
* Должна быть возможность запускать проект с конфигами в ```.toml```
```
api -format .env -path configs/.env
```
* Должна быть возможность запускать проект с конфигами в ```.env```
```
api -format .toml -path configs/api.toml
```
* Должна быть возможность запускать проект с дефолтными параметрами (дефолтным будем считать ```apiserver.toml```,
если его нет, то запускаем с значениями из структуры ```Config```)
*/

func init() {
	//Скажем, что наше приложение будет на этапе запуска получать путь до конфиг файла из внешнего мира
	flag.StringVar(&typeConfigFile, "format", ".toml", "format config file in .toml format or .env")
	flag.StringVar(&configPath, "path", "configs/apiserver.toml", "path to config file in .toml format")

}

func main() {
	//В этот момент происходит инициализация переменной configPath значением
	flag.Parse()
	log.Println("typeConfigFile (after Parse2): ", typeConfigFile, "configPath (after Parse2): ", configPath)
	//server instance initialization
	config := apiserver.NewConfig()
	switch typeConfigFile {
	case ".toml":
		_, err := toml.DecodeFile(configPath, config) // Десериалзиуете содержимое .toml файла
		if err != nil {
			log.Printf("can not find configs file. using default file: %s Error: %s", configPathDefault, err)
			_, err = toml.DecodeFile("configs/apiserver.toml", config)
			if err != nil {
				log.Fatal("can not find default configs file: ", err)
			}
		}
	case ".env":
		err := godotenv.Load(configPath)
		if err != nil {
			log.Printf("can not find configs file. using default file: %s Error: %s", configPathDefault, err)
			_, err = toml.DecodeFile("configs/apiserver.toml", config)
			if err != nil {
				log.Fatal("can not find default configs file: ", err)
			}
		} else {

			config.BindAddr = os.Getenv("bind_add")
			config.LogLevel = os.Getenv("log_level")
			config.Store.DatabaseURL = os.Getenv("database_url")
		}

	default:
		{
			log.Printf("can not find configs file. using default file: %s", configPathDefault)
			_, err := toml.DecodeFile("configs/apiserver.toml", config)
			if err != nil {
				log.Fatal("can not find default configs file: ", err)
			}
		}

	}
	log.Println("config:", config)

	//server instance
	s := apiserver.New(config)

	//server start
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
