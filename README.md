# Распределенный вычислитель арифметических выражений

<details><summary><b>Задание</b></summary>

Пользователь хочет считать арифметические выражения.
Он вводит строку `2 + 2 * 2` и хочет получить в ответ `6`.
Но наши операции сложения и умножения (также деления и вычитания) выполняются **"очень-очень" долго**.
Поэтому вариант, при котором пользователь делает http-запрос и получает в качетсве ответа результат, **невозможна**.
Более того: вычисление каждой такой операции в нашей **"альтернативной реальности"** занимает **"гигантские"** вычислительные мощности.
Соответственно, каждое действие мы должны уметь выполнять отдельно и масштабировать эту систему можем добавлением вычислительных мощностей в нашу систему в виде новых "**машин**".
Поэтому пользователь, присылая выражение, получает в ответ идентификатор выражения и может с какой-то периодичностью уточнять у сервера "не посчиталость ли выражение"?
Если выражение наконец будет вычислено - то он получит результат.
Помните, что некоторые части арфиметического выражения можно вычислять **параллельно**.

## Back-end часть

### Состоит из 2 элементов:

- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
- Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.

### Оркестратор
Сервер, который имеет следующие endpoint-ы:

- Добавление вычисления арифметического выражения.
- Получение значения выражения по его идентификатору.
- Получение списка всех выражений.
- Получение задачи для выполнения.
- Приём результата обработки данных.


### Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения. При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды.
</details>

## Инструкция к запуску:
1. Склонировать проект или скачать `git clone github.com/katenester/Distributed_computing`
2. Установить все зависимости из go.mod `go mod download` `go mod tidy`
3. Переходим в папку с проектом на компьютере. `cd Distributed_computing`
4. Открываем два терминала (win+R->cmd). В первом окне ввести ` go run cmd/orchestrator/main.go` . Во втором окне ввести `go run cmd/agent/main.go `. Должны запустится процессы и идти логи. Для отключения процессов Ctrl+s

## Тестирование:
<details><summary><b>Примеры использования с помощью curl</b></summary>
</details>
<details><summary><b>Postman</b></summary>

</details>