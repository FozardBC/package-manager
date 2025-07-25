# 📦 `pm` — Простой пакетный менеджер на Go

`pm` — это легковесный CLI-инструмент для упаковки, загрузки и установки пакетов по SSH. Написан на Go. Подходит для развёртывания конфигураций, дистрибуции бинарников или управления внутренними пакетами.

---

##  Возможности

- Упаковка файлов в `.tar.gz` по маскам
- Загрузка и скачивание пакетов через **SSH/SFTP**
- Поддержка **версионности**: `>=1.10`, `<=2.0`, точные версии
- Конфигурация в формате **JSON или YAML**
- Исключение файлов по шаблону (например, `*.tmp`)
- Поддержка **SSH-ключей** 
- Один бинарный файл — нет зависимостей


- ## 📥 Установка

### через Makefile


```bash
git clone https://github.com/FozardBC/package-manager
cd pm

make          # Показать справку (по умолчанию)
make build    # Собрать бинарник: pm + Добавить бинарник в /usr/bin
make mod      # Обновить зависимости: go mod tidy
make clean    # Удалить бинарник

```

- # Ручками
```bash
# Клонируй репозиторий
git clone https://github.com/FozardBC/package-manager
cd pm

# Создайте директорию для SSH key
mkdir vault
# Поместите в vault ключ

# Определите переменные окружения
# (опционально) Создайте .env и определите
# (опционально)
export SSH_HOST=localhost
export SSH_USER=root
export SSH_PORT=2222
export SSH_KEY_PATH=./vault/test-key/test-key

# Соберите бинарник
go build -o pm

# (Опционально) Добавь в PATH
sudo mv pm /usr/local/bin/
```

🛠️ Использование
Создаёт архив из файлов по конфигу и загружает его на сервер.
```
pm create <файл-конфигурации>
```
Скачивает и устанавливает пакеты, указанные в конфиге.

```
pm update <файл-списка>
```

📁 Формат конфигурации
packet.json — Какие файлы упаковать
```json
{
  "name": "myapp",
  "ver": "1.10", # версия в формате x.xx
  "targets": [
    "./configs/*.yaml",
    {
      "path": "./bin/*", 
      "exclude": ["*.tmp", "dev-*"] # исключения
    }
  ]
}
```

packages.json — Какие пакеты установить
```json
{
  "packages": [
    { "name": "myapp", "ver": ">=1.10" },
    { "name": "utils", "ver": "1.5" },
    { "name": "logger", "ver": "<=2.0" }
  ]
}
```

🔐 Аутентификация
По SSH-ключу
Настройки подключения (хост, пользователь, путь к ключу) можно задавать в коде или через переменные окружения.


## Для тестирования приложения через тестовы SSH-сервер в Docker
```bash
docker run -d --name ssh-test -p 2222:22 rastasheep/ubuntu-sshd:18.04

#копирование ключа в контейнер
docker cp test_key.pub ssh-test-server:/tmp/id_rsa.pub

docker exec -it ssh-test-server /bin/bash

mkdir -p /root/.ssh
cat /tmp/id_rsa.pub >> /root/.ssh/authorized_keys
chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys
exit
```

```
EXPORT PM_SSH_HOST=localhost
EXPORT PM_SSH_USER=root
EXPORT PM_SSH_PORT=2222
EXPORT PM_SSH_KEY_PATH=./vault/#ssh-key

#или создать файл .env в корневой папке

go build -o .

sudo mv pm /usr/local/bin/

pm ..
```

