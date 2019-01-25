#!/bin/sh

xterm -e bash -c 'vgo run main.go -port 3000' &
xterm -e bash -c 'vgo run main.go -port 3001 -peers tcp://localhost:3000' &
xterm -e bash -c 'vgo run main.go -port 3002 -peers tcp://localhost:3000'