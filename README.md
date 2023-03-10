# Задание
## MemcLoad v2

Задание
: нужно переписать Python версию memc_load.py на Go. Программа по-прежнему парсит и заливает в мемкеш поминутную выгрузку логов трекера установленных приложений. Ключом является тип и идентификатор устройствачерез двоеточие, значением являет protobuf сообщение(
https://github.com/golang/protobuf
).

```
$ ls -lh /Users/s.stupnikov/Coding/otus/*.tsv.gz
  -rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:09 /Users/s.stupnikov/Codin
  -rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:17 /Users/s.stupnikov/Codin
  -rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:25 /Users/s.stupnikov/Codin$ 
```
gunzip -c /Users/s.stupnikov/Coding/otus/20170929000000.tsv.gz | head -3idfae7e1a50c0ec2747ca56cd9e1558c0d7c67.7835424444-22.804400547idfaf5ae5fe6122bb20d08ff2c2ec43fb4c4-104.68583244-51.24448376gaid3261cf44cbe6a00839c574336fdf49f6137.79083956756.8403675248
Ссылки на tsv.gz файлы:
```
https://cloud.mail.ru/public/2hZL/Ko9s8R9TA
https://cloud.mail.ru/public/DzSX/oj8RxGX1A
https://cloud.mail.ru/public/LoDo/SfsPEzoGc
```
Важно обрабатывать файлики в хронологическом порядке. В данном случае подэтим имеется в вижу, что после обработки нужно переименовывать файл,префиксировав имя точкой, последовательно и хронологически. Заливать жеможно параллельно.
# Цель задания: поработать с моделью конкурентности отличной от Python, получить навык решения задач на новом языке.


Необходимые библиотеки:
- go get gotest.tools/v3
- go get -u github.com/google/go-cmp/cmp
- go get github.com/stretchr/testify/assert
- go get github.com/stretchr/testify

# Установка и запуск memcached
https://commaster.net/posts/installing-memcached-windows/

https://go.dev/doc/code#Testing
Запуск тестов (в папке go_daemon_protobuf):
go test

# Запуск скрипта:
```
go run .\memc_load.go --pattern="C:\memcached\sample.tsv.gz" --idfa="127.0.0.1:11211" --gaid="127.0.0.1:11211" --adid="127.0.0.1:11211" --dvid="127.0.0.1:11211"
```

# Убедится в том, что данных сохранились:
```
telnet 127.0.0.1 11211
get idfa:1rfw452y52g2gq4g
```
где idfa:1rfw452y52g2gq4g
- это ключ (его можно взять в файле или в логах)
----------------------
Команда генерации кода на основе proto-файла
- python:
https://protobuf.dev/reference/python/python-generated/
```
protoc --proto_path="E:\python scripts\go_daemon" --python_out="E:\python scripts\go_daemon\task" "E:\python scripts\go_daemon\appsinstalled.proto"
```
- go
https://protobuf.dev/reference/go/go-generated/
```
protoc --go_opt=Mappsinstalled.proto=example.com/project/protos/fizz --go_out=. --go_opt=paths=source_relative appsinstalled.proto 
```

----------------------
# Замечания:
1) Я вынес тесты в отдельную функцию
2) я не нашёл аналога unpacked.ParseFromString(packed)
и решил сравнивать с ожидаемым результатом
3) insert_appsinstalled значение по умолчанию не делал
4) https://github.com/bradfitz/gomemcache/
слишком старая библиотека
5) приходится постоянно переключаться между 
```
go env -w GO111MODULE=off
```
и 
```
go env -w GO111MODULE=auto
```
```
go mod vendor
go: modules disabled by GO111MODULE=off; see 'go help modules'
...
 go mod vendor
 go build
memc_load.go:8:2: package go_daemon_protobuf is not in GOROOT (C:\Program Files\Go\src\go_daemon_protobuf)
go env -w GO111MODULE=off
```
```
go run .\memc_load.go -log -dry
go run .\memc_load.go -dvid="127.0.0.1:444"
go run .\memc_load.go --pattern="E:\python scripts\go_daemon\test_data\sample.tsv.gz"
```


