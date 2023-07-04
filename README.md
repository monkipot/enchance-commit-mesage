# Enchance your commit messages

A simple piece of code to force your team to write relevant and conventional commit messages.
This code displays commit messages (from here https://whatthecommit.com/index.txt ). If your message is conventional,
then it will display a message, otherwise it will turn your message into a random message. :)

## Setup

```bash
echo -e "
function git() {
  if [[ \$1 == \"commit\" ]]; then
    command git commit \"\$@\" && ~/path/to/enhance-commit-messages
  else
    command git \"\$@\"
  fi
}
" >> ~/.bashrc

source ~/.bashrc
```
