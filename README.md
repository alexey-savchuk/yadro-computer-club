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

### Файлы для тестовых запусков
- [./testdata/01_default.txt](testdata/01_default.txt) - пример из задания
- [./testdata/02_invalid_event_time.txt](testdata/02_invalid_event_time.txt) -
  некорректное время события: ```09:63 3 client1```
- [./testdata/03_invalid_event_id.txt](testdata/03_invalid_event_id.txt) -
  некорректный идентификатор события: ```09:41 5 client1```
- [./testdata/04_invalid_event_body.txt](testdata/04_invalid_event_body.txt) -
  некорректное тело события: ```09:54 2 client1```

Пример запуска
```sh
make build
make run FILE="testdata/04_invalid_event_body.txt"
```
получаем сообщение об ошибке, и программа завершается
```
failed to parse event: invalid event "09:54 2 client1" format, event #2 requires <name> <table number> body
make: *** [Makefile:7: run] Error 1
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
