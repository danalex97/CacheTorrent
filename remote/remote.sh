ARGS=$@
if [ "$1" -eq '-h' ];
then
  python3 remote.py -h
elif [ "$1" -eq '--help' ];
then
  python3 remote.py -h
else
  echo "Running job: $ARGS"
  nohup python3 remote.py $ARGS > /dev/null 2>&1 &
fi
