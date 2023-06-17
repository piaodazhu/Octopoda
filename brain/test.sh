#!/bin/bash

ADDR=/root/Octopoda/brain
(sleep 3 && echo "ok" >> $ADDR/foobar &)