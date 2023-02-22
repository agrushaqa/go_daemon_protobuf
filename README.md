# Задание
## MemcLoad v2

Задание
: нужно переписать Python версию memc_load.py на Go. Программа по-прежнему парсит и заливает в мемкеш поминутную выгрузку логов трекера установленных приложений. Ключом является тип и идентификатор устройствачерез двоеточие, значением являет protobuf сообщение(
https://github.com/golang/protobuf
).
$ ls -lh /Users/s.stupnikov/Coding/otus/*.tsv.gz

-rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:09 /Users/s.stupnikov/Codin
-rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:17 /Users/s.stupnikov/Codin
-rw-r--r-- 1 s.stupnikov staff 506M 29 сен 12:25 /Users/s.stupnikov/Codin$ 

gunzip -c /Users/s.stupnikov/Coding/otus/20170929000000.tsv.gz | head -3idfae7e1a50c0ec2747ca56cd9e1558c0d7c67.7835424444-22.804400547idfaf5ae5fe6122bb20d08ff2c2ec43fb4c4-104.68583244-51.24448376gaid3261cf44cbe6a00839c574336fdf49f6137.79083956756.8403675248
Ссылки на tsv.gz файлы:
https://cloud.mail.ru/public/2hZL/Ko9s8R9TA
https://cloud.mail.ru/public/DzSX/oj8RxGX1A
https://cloud.mail.ru/public/LoDo/SfsPEzoGc
Важно обрабатывать файлики в хронологическом порядке. В данном случае подэтим имеется в вижу, что после обработки нужно переименовывать файл,префиксировав имя точкой, последовательно и хронологически. Заливать жеможно параллельно.
Цель задания
: поработать с моделью конкурентности отличной от Python, получитьнавык решения задач на новом языке.

# requirements.txt
## create
pip freeze > requirements.txt
pip install -r requirements.txt
## use
# code style
## isort
python -m pip install isort
### run 
isort .
## mypy
python -m pip install mypy
### run 
mypy .
## flake8
python -m pip install flake8
### run
flake8 --exclude venv,docs --ignore=F401
## code coverage
pip install coverage
### run
coverage run C:\Users\agrusha\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.10_qbz5n2kfra8p0\LocalCache\local-packages\Python310\site-packages\behave\__main__.py
в файле .coveragerc нужно указать исходники

# Pytest - Run Tests in Parallel
## install
```pip install pytest-xdist```

also:
https://pypi.org/project/pytest-parallel/
## run
```pytest -n 2 test_common.py```
