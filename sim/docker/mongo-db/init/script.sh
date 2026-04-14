#!/usr/bin/env bash
mongo admin -u root -p root --eval "db.createUser({user: 'duel', pwd: 'master', roles: [{role: 'dbOwner', db: 'duel-masters'}]});"
