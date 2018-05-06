ARGS=$@

function on_done() {
  echo "Job dispached."
  exit 0
}

if [ "$1" = '-h' ]
then
  python3 remote/remote.py -h
elif [ "$1" = '--help' ]
then
  python3 remote/remote.py -h
else
  pid=$$
  echo "Dispaching job: $ARGS"
  echo "PID: $pid"
  nohup python3 remote/remote.py -notify=$pid $ARGS > /dev/null 2>&1 &
  # python3 remote/remote.py -notify=$pid $ARGS &

  # Wait for SIGUR1 signal
  trap "on_done" SIGUSR1
  while :; do
    sleep 1
  done
fi
