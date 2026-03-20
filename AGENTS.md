# Project Instructions

## After code editing
* When you add, remove or modify files under pkg/formatter/templates/text/, run `make bindata` to regenerate pkg/formatter/templates/bindata.go. Commit updated bindata.go together with the template changes.
* When you add, remove of modify commands, run `make docs` to regenerate the files under docs/. Commit updated docs together with the command changes.
* Run `make format` and `make lint` to verify your code style. If these result in changes, evaluate if the changes are sensible. If they are sensible, commit them together with your code changes.
* Run `make unit` to run the unit tests to verify your code changes.
