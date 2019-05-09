#!/bin/sh
set -eu

if [ "$#" -lt 1 ]; then
  echo "usage: mirror.sh <github_ssh_clone_url>" 1>&2
fi

GITHUB_MIRROR_URL="$1"

# determine current commit
GIT_COMMIT="$(git rev-parse HEAD)"

ensure_ssh() {
  # ensure our SSH allows us to talk to GitHub.
  # Even if we ended up talking to a malicious host, all we'd do is push our
  # open source code to them? Not really concerned about that.
  KNOWN_HOSTS_FILE="$HOME/.ssh/known_hosts"
  touch "$KNOWN_HOSTS_FILE"
  chmod 0600 "$KNOWN_HOSTS_FILE"
  ssh-keyscan github.com > "$KNOWN_HOSTS_FILE"
}

# if our current commit is on the master branch, push to GitHub
if git branch --format '%(refname:lstrip=2)' --contains "$GIT_COMMIT" | grep 'master'; then
  ensure_ssh
  git remote add origin "$GITHUB_MIRROR_URL" && \
  git checkout master && \
  git push -u origin master
fi
