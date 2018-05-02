ARGS=$@
ssh ad5915@shell3.doc.ic.ac.uk << ENDSSH
  echo "Connected to shell."
  echo "Job: $ARGS"
  ssh -f -o StrictHostKeyChecking=no ad5915@matrix01.doc.ic.ac.uk "
    echo 'Dispaching: $ARGS'
    cd ~/golang/src/github.com/danalex97/nfsTorrent
    nohup python3 remote.py $ARGS > /dev/null 2>&1 &
    echo 'Job dispached.'
    exit
  " &
ENDSSH &
