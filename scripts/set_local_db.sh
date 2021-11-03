#!/bin/bash

pg_ctl -D ./data/toc_trader -l ./data/toc_trader/logfile stop
rm -rf ./data
mkdir -p ./data/toc_trader

initdb ./data/toc_trader
pg_ctl -D ./data/toc_trader -l ./data/toc_trader/logfile start

echo "\du
CREATE ROLE postgres WITH LOGIN PASSWORD 'asdf0000';
ALTER USER postgres WITH SUPERUSER;
\du" > sql_script

psql postgres -f sql_script
rm -rf sql_script

pg_ctl -D ./data/toc_trader -l ./data/toc_trader/logfile stop
