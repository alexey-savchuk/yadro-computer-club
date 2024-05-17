### Формулировка задания
Можно посмотреть [здесь](TASK.md)

### Требования
- Docker
- Makefile
- realpath[^1]

### Как запустить
```sh
make build
make run FILE="/path/to/file"
```
или
```sh
make build run FILE="/path/to/file"
```

[^1]: realpath - стандартная команда Linux для получения абсолютного пути.
    Скорее всего уже установлена по-умолчанию, но на всякий случай дам ссылки
    - [realpath в Ubuntu](https://manpages.ubuntu.com/manpages/trusty/en/man1/realpath.1.html)
    - [realpath в Arch](https://man.archlinux.org/man/realpath.1.en)
