# !/bin/bash

# . /kb/deployment/user-env.sh

#python3 ./scripts/prepare_deploy_cfg.py ./deploy.cfg ./work/config.properties

if [ -f ./work/token ] ; then
  export KB_AUTH_TOKEN=$(<./work/token)
fi

if [ $# -eq 0 ] ; then
  sh ./scripts/start-server.sh
elif [ "${1}" = "test" ] ; then
  sh ./scripts/run-tests.sh
elif [ "${1}" = "async" ] ; then
  sh ./scripts/run-async.sh
elif [ "${1}" = "init" ] ; then
  echo "Initialize module"
elif [ "${1}" = "bash" ] ; then
  bash
elif [ "${1}" = "report" ] ; then
  echo "Copying compilation report to /work"
  cp ./scripts/compile_report.json ./work
else
  @echo "unknown entrypoint command"
fi
