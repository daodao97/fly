#!/bin/bash
set -x

rm -rf .git
git init
git config user.name daodao97 
git config user.email daodao97@foxmail.com
git add .
git commit -m '·'
git remote add origin git@github.com:daodao97/fly.git
git push origin master -f

