#!/bin/bash

# The Docker entrypoint runs this script automatically during initialization.
# We connect as root using the password provided in the environment variables
# and grant read access to the specific schema for the auto-created user.

mysql -u root <<-EOSQL
    GRANT SUPER ON *.* TO '$MYSQL_USER'@'%';
    GRANT SELECT ON performance_schema.* TO '$MYSQL_USER'@'%';
    FLUSH PRIVILEGES;
EOSQL
