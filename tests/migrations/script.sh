#!/bin/bash

draconctl apply --url "postgres://postgres:postgres@localhost:5432/?sslmode=disable" /etc/migrations
draconctl revert --url "postgres://postgres:postgres@localhost:5432/?sslmode=disable" /etc/migrations
