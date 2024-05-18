# Тестовое задание для стажировки "Импульс" YADRO 2024 // GO

## Инструкция по запуску

---
Запуск контейнера с программой:

Сборка:

```zsh
git clone https://github.com/GerogeGol/yadro-test-problem 
cd yadro-test-problem
docker build --tag yadro-problem .

```

Запуск:

```zsh
docker run yadro-problem:latest /main tests/basic.txt
```

Также можно посмотреть результаты работы программы на тестовых входных данных:

```zsh
docker run -it yadro-problem:latest bash
cd /app/tests/out
cat basic.out
...
```
