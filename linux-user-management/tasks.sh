#!/usr/bin/env bash
# The three operations from the question, run as ec2-user with sudo.
#
# Each command is printed after a shell prompt and then run, so the output
# reads like the terminal session it is.
#
# usage: tasks.sh [0|1|2|3|all]      (0 shows the starting state)
#        tasks.sh 2a|2b|3a|3b        (the same steps in two halves)
set -uo pipefail

PROJECT=/opt/Hdfc_Bank-ClientApp

run() {
  local dir=${PWD/#$HOME/\~}
  printf '\n[%s@%s %s]$ %s\n' "$(whoami)" "${HOSTNAME%%.*}" "${dir##*/}" "$*"
  eval "$*"
}

# Starting state: the machine, both groups, every home directory, the project.
task0() {
  echo "===== Starting state ====="
  run "grep PRETTY_NAME /etc/os-release"
  run "getent group Developer Tester"
  run "ls /home"
  run "ls $PROJECT"
}

# (i) A new user, Developer4, joins and must be in the Developer group.
task1() {
  echo "===== (i) Add Developer4 to the Developer group ====="
  run "getent group Developer"
  run "sudo useradd -m Developer4"
  run "sudo usermod -aG Developer Developer4"
  run "id Developer4"
  run "getent group Developer"
  run "ls -ld /home/Developer4"
}

# (ii) Developer3 leaves: remove the account and its home directory completely.
task2a() {
  echo "===== (ii) Remove Developer3 completely ====="
  run "id Developer3"
  run "sudo find /home/Developer3"
  run "sudo userdel -r Developer3"
}
task2b() {
  echo "===== (ii) Check that nothing of Developer3 is left ====="
  run "id Developer3"
  run "sudo grep -s Developer3 /etc/passwd /etc/shadow /etc/group /etc/gshadow || echo 'no entries left'"
  run "ls /home"
  run "getent group Developer"
  run "sudo find / -xdev -nouser 2>/dev/null | wc -l"   # files whose owner no longer exists
}

# (iii) Show what is inside Hdfc-Pack.zip, then delete a file from it.
task3a() {
  echo "===== (iii) Show the contents of Hdfc-Pack.zip ====="
  cd "$PROJECT" || exit 1
  run "ls -l"
  run "unzip -l Hdfc-Pack.zip"
}
task3b() {
  echo "===== (iii) Delete a file from Hdfc-Pack.zip ====="
  cd "$PROJECT" || exit 1
  run "sudo zip -d Hdfc-Pack.zip Hdfc-Pack/logs/debug.log"
  run "unzip -l Hdfc-Pack.zip"
}

case "${1:-all}" in
  0) task0 ;;
  1) task1 ;;
  2) task2a; task2b ;;
  2a) task2a ;;
  2b) task2b ;;
  3) task3a; task3b ;;
  3a) task3a ;;
  3b) task3b ;;
  all) task0; task1; task2a; task2b; task3a; task3b ;;
  *) echo "usage: $0 [0|1|2|3|all|2a|2b|3a|3b]" >&2; exit 2 ;;
esac
