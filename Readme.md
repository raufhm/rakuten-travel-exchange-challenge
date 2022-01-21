# Documentation

## 1. Clone or extract file

extract this zip file into your golang development environment

## 2. Installation

setup your local machine by install 

- Docker 

- Database Migration CLI (https://github.com/golang-migrate/migrate)

- gin-gonic

- lib/pq

## 3. Run Command

once everyhing is installed please execute below command. for futher detail, please go to ```Makefile``` in root folder. it will contain all command that will be used to use the application.

- ```sudo make postgres```
  this command is to setup database in docker with name postgres12.
  make sure the images is is created and container is up with command ```sudo docker images``` and ```sudo docker ps -a```

- ```sudo make createdb```
  this command is to create database

- ```sudo make dropdb```
  this command is to drop database

- ```sudo make migrateup```
  this command is to apply data migration. detail in ./db/migration/000001_init_schema.up.sql

- ```sudo make migratedown```
  this command is to delete our table that we've created. detail in ./db/migration/000001_init_schema.down.sql

- ```sudo make gotodb```
  this command is to go our database. we may check whether the table is created by excute ```\dt``` to find table list

## 4. Run Server

- execute in terminal ```go run main.go``` or simply hit f5 in to run the server. server will be run on localhost:8080

- program will auto load information since start running from xml url (https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml) and bulk import all data into our database. 

  to avoid duplication, please clear database is cleared before re-run the server.

- from your browser or postman, please run below endpoint

    - http://localhost:8080/rates/latest to get the latest rate.

        - to see the response, please find in ./response/latestRate.json

    - http://localhost:8080/rates/2021-11-10 to get rate based on date request

        - to see the response, please find in ./response/rateByDate.json 

    - http://localhost:8080/rates/analyze to get rate by doing calculation base on loaded data from xml by return maximum, minimum, and average rate.

        - to see the response, please find in ./response/rateAnalyze.json
