ARGS=$@
echo "Running job: $ARGS"
nohup python3 remote.py $ARGS > /dev/null 2>&1 &
