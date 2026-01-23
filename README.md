# Test Assignment
Тестовое задание для itk academy

Для запуска приложения требуется рядом с main.go создать файл config.env на основе config.env.example(он находится в /cmd/wallet/)

Для запуска docker-compose:
```bash
make up
```

Для запуска приложения локально(также требуется config.env):
```bash
go run cmd/wallet/main.go
```

Для локального запуска потребуется провести автогенерацию файлов:
```bash
make oapi-code-gen
```
Либо запустить:
```bash
make all
```

Для проверки тестов, нужно ввести команду:
```bash
make tests
```

Все команды описаны в Makefile
