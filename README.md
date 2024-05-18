### Формулировка задания
Можно посмотреть [здесь](TASK.md)

### Требования
- Docker
- Makefile
- realpath[^1]

### Как запустить
```sh
make build                    # сборка
make run FILE="/path/to/file" # запуск
```
или
```sh
make build run FILE="/path/to/file"
```

### Дополнительно
Удалить Docker-образ, созданный командой ```make build```, можно так
```sh
make clean
```


[^1]: realpath - стандартная команда Linux для получения абсолютного пути.
    Скорее всего уже установлена по-умолчанию, но на всякий случай дам ссылки
    - [realpath в Ubuntu](https://manpages.ubuntu.com/manpages/trusty/en/man1/realpath.1.html)
    - [realpath в Arch](https://man.archlinux.org/man/realpath.1.en)
