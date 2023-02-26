# Задание
## MemcLoad v2

Задание
: нужно переписать Python версию memc_load.py на Go. Программа по-прежнему парсит и заливает в мемкеш поминутную выгрузку логов трекера установленных приложений. Ключом является тип и идентификатор устройствачерез двоеточие, значением являет protobuf сообщение(
https://github.com/golang/protobuf
).
$ ls -lh /Users/s.stupnikov/Coding/otus/*.tsv.gz
```
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
# Цель задания: поработать с моделью конкурентности отличной от Python, получитьнавык решения задач на новом языке.


Необходимые библиотеки:
- go get gotest.tools/v3
- go get -u github.com/google/go-cmp/cmp
- go get github.com/stretchr/testify/assert
- go get github.com/stretchr/testify


https://go.dev/doc/code#Testing
Запуск тестов (в папке go_daemon_protobuf):
go test

----------------------
1) Я вынес тесты в отдельную функцию
2) я не нашёл аналога unpacked.ParseFromString(packed)
и решил сравнивать с ожидаемым результатом
3) insert_appsinstalled значение по умолчанию не делал
4) https://github.com/bradfitz/gomemcache/
слишком старая библиотека
5) приходится постоянно переключаться между 
go env -w GO111MODULE=off
и 
go env -w GO111MODULE=auto
PS E:\python scripts\go_daemon> go mod vendor
go: modules disabled by GO111MODULE=off; see 'go help modules'
...
PS E:\python scripts\go_daemon> go mod vendor
PS E:\python scripts\go_daemon> go build
memc_load.go:8:2: package go_daemon_protobuf is not in GOROOT (C:\Program Files\Go\src\go_daemon_protobuf)
PS E:\python scripts\go_daemon> go env -w GO111MODULE=off

go run .\memc_load.go -log -dry
go run .\memc_load.go -dvid="127.0.0.1:444"
go run .\memc_load.go --pattern="E:\python scripts\go_daemon\test_data\sample.tsv.gz"

6) сейчас вызов
go run .\memc_load.go --pattern="E:\python scripts\go_daemon\test_data\sample.tsv.gz"

приводит к ошибке
2023/02/26 15:24:47 Cannot write to memc 127.0.0.1:33013: memcache: connect timeout to 127.0.0.1:33013
exit status 1
