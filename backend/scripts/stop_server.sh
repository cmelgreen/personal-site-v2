#!/bin/bash

if pidof server; then
    killall server
fi