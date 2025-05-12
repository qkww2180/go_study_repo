#!/bin/bash

i=0
while [ true ]
do
  let i++
  curl localhost/blog
done