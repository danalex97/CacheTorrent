if [ -z ${BROWSER+x} ]
then
  BROWSER=firefox
fi

HERE=`pwd`
# Run backend server
PYTHONPATH=$HERE python3 visual/run.py

# Run frontend display
$BROWSER visual/index.html
