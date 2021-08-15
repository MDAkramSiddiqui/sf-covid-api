#!/bin/bash

export ENV=DEVELOPMENT
export MONGO_DB_URL=mongodb+srv://sf-covid-admin:{{password}}@cluster0.hkjao.mongodb.net/{{db_name}}?retryWrites=true&w=majority
export MONGO_DB_PASSWORD=Akram077
export DB_NAME=covid