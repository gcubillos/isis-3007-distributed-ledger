#!/bin/sh

vgo run main.go -port 3002 -peers tcp://192.168.0.18:3000 tcp://192.168.0.18:3001
