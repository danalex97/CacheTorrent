ARGS=$@
echo "Running job: $ARGS"
nohup python3 reomte.py $ARGS > /dev/null 2>&1 &
