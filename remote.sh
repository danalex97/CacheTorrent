ARGS=$@
ssh ad5915@shell1.doc.ic.ac.uk << ENDSSH
  echo "Connected to shell1."
  echo "Job: $ARGS"
  ssh -tt -o StrictHostKeyChecking=no ad5915@matrix01.doc.ic.ac.uk "
    echo 'Remote job dispached: $ARGS'
    cd ~/golang/src/github.com/danalex97/nfsTorrent
    python3 remote.py $@
    exit
  "
  exit
ENDSSH
