# common-go-modules
A repo to store common functionality for the go services like logging, error handling, rest client, etc

#### Setup Commit Hooks

1. Make sure the PATH var in your .bashrc file has this
```bash
	$HOME/bin
```
2. Install pre-commit application to allow pre commit hooks to run correctly:
```bash
    $ curl https://pre-commit.com/install-local.py | python -
```
3. Link pre-commit to your local repo. Run the below command from the root of your repo
```bash
    $ pre-commit install
    pre-commit installed at .git/hooks/pre-commit