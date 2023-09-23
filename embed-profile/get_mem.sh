#!/usr/bin/env bash

ps x -o rss,vsz,command $1 | awk 'NR>1 {$1=int($1/1024)"M"; $2=int($2/1024)"M";}{ print ;}'
