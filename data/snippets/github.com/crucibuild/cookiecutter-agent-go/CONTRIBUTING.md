# General
In order to work effectively as a decentralized team, we are doing small and quick iterations on branches followed by a Pull Request when deemed stable locally.

# Rules
- At least a reviewer needs to be assigned for each Pull Requests.
- A Pull Request should not consist of more than 100 lines of diff as
  it will be difficult to review. However, for a few limited set of
changes (like refactoring), we can bypass this rule.
- Except for user created issues, an issue should be a task evaluated to
  2-3 days of work.
- A PR should always have an associated issue.
- A branch associated to an issue should follow this naming scheme:
```
<TYPE>-<NUMBER>-<TITLE>
```
For example:
```
 bug-3-Makefile-missing-examples
```
- Commit message should be formatted this way:

```
[COMPONENT] #<ISSUE> <MESSAGE>

<DESCRIPTION>
```

COMPONENT is a facultative field.
