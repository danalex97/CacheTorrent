ARGS=$@
if [ "$1" = '-h' ]
then
  python3 remote/remote.py -h
elif [ "$1" = '--help' ]
then
  python3 remote/remote.py -h
else
  echo "Running job: $ARGS"
  nohup python3 remote/remote.py $ARGS > /dev/null 2>&1 &
fi
