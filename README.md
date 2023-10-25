# gitlab-ci-mr
Возможности при создании ветки в CI:

- Пуш `feature/XXX-123` создает ветку `review/XXX-123` и MR к ним
- Пуш `review/XXX-123` или `hotfix/XXX-123` создает MR в `main` (или в другую стандартную ветку из проекта)

В качестве названия MR использует значение из `CI_COMMIT_MESSAGE`

Существенный минус: токен `PAT` (personal access token) будут видеть все.

Лучше использовать [Git Hooks](#git-hooks)

### Использование

```yaml
# .gitlab-ci.yml
image: alpine:latest

stages:
  - create-mr

create-mr:
  stage: create-mr
  script:
    - apk --no-cache add ca-certificates
    - ./gitlab-ci-mr
  only:
    - /^(feature|review|hotfix).*/
```

### Параметры

Необходимо задать 1 переменную `PAT` (personal access token). Взять из [профиля](https://gitlab.com/-/profile/personal_access_tokens).
`CI_JOB_TOKEN` не подходит, нет прав.

1. https://gitlab.com/gitlab-org/gitlab/-/issues/17511
2. https://docs.gitlab.com/ee/ci/jobs/ci_job_token.html

Остальные - это стандартные переменные https://docs.gitlab.com/ee/ci/variables/predefined_variables.html

## Другие способы

- [Git Hooks](#git-hooks)
- [Push options](#push-options)
- ???

### Git Hooks

В git нет хука `post-push`, но есть костыль.

Тут [инструкция](./with-git-hooks/README.md)

### Push options

Минусы
- Не создает ветки

В [GitLab](https://docs.gitlab.com/ee/user/project/push_options.html) можно использовать параметры для создания MR при выполнении `git push`.

```shell
# пример создания alias для main ветки
git config --global alias.mwps "push -o merge_request.create -o merge_request.target=main"
git mwps origin <local-branch-name>
```

Если указать `merge_request.target=<несуществующая-ветка>` - ошибка
```
WARNINGS: Error encountered with push options
'merge_request.create' 'merge_request.target=review/not-exists': Target
branch group/repo:review/not-exists does not exist
```
