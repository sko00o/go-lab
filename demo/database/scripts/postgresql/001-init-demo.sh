#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  create role root with login password 'toor';
  create database demo_test with owner root;
EOSQL