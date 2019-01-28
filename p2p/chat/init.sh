#!/bin/sh

xterm -e bash -c "vgo run main.go -port 3000" &

for i in {1..2}
do
   xterm -e bash -c "vgo run main.go -port 300$i -peers tcp://localhost:3000" &
done
