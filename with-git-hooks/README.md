# Git Hook

Используем `bash` скрипты для создания веток и MR.

#### Что именно делает скрипт

Сначала пользователь пушит новую ветку `feature/TASK-1010`, затем скрипт:
* проверяет наличие переменной окружения `MY_ZQSHDZW_PP`
  * если не существует: выход из скрипта
* получает ID проекта в GitLab из значения `git remote get-url origin`
* получает ID текущего пользователя
* сохраняет и не перезапрашивает эти переменные в файл `.githooks/.env`
* проверяет создана ли ветка `review/TASK-1010`
  * если не существует: создает, за основу берет `main`
* проверяет существует ли открытый Merge Request из `feature/TASK-1010` в `review/TASK-1010`
  * если не существует: создает с данными
    * название: текст коммита
    * назначает текущего юзера
* выдает ссылку на Merge Request

### Установить зависимости
`# sudo pacman -S curl jq git`

## Установка

1. Скопировать содержимое `.githooks` в репозиторий проекта
2. (Если нужно) Заменить 2 значения в `.githooks/post-push`
```shell
MAIN_BRANCH=main # или master
GITLAB_URL=gitlab.com # или gitlab.my-company.org
```
3. Добавить `.githooks/.env` в `.gitignore`
4. Установить локально переменную окружения `MY_ZQSHDZW_PP`
```shell
# https://gitlab.com/-/profile/personal_access_tokens
# echo 'export MY_ZQSHDZW_PP="gpat-"' >> ~/.zshrc
```
5. Алиас для замены `/.git/hooks` на `/.githooks` для хранения в Git.
```shell
git config --local core.hooksPath .githooks/
```
6. Заливаем изменения `# git commit -m "[TASK-1010] add git hook"` & `git push`
7. Видим сообщение `MR [ feature/TASK-1010 -> review/TASK-1010 ] https://gitlab.com/:repo/-/merge_requests/1`
