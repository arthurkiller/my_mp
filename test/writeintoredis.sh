#!/bin/bash

# set for the userinfo
redis-cli -p 16380 HMSET a1 Uid a1 Name a1 Pic "" FansNum 2 FollowNum 2 MeipaiNum 0
redis-cli -p 16380 HMSET b2 Uid b2 Name b2 Pic "" FansNum 2 FollowNum 2 MeipaiNum 0
redis-cli -p 16380 HMSET c3 Uid c3 Name c3 Pic "" FansNum 2 FollowNum 2 MeipaiNum 0

#set for the userfans
redis-cli -p 16385 SADD a1 b2 c3
redis-cli -p 16385 SADD b2 c3 a1 
redis-cli -p 16385 SADD c3 a1 b2 

#set fot thr userfollow
redis-cli -p 16386 SADD a1 b2 c3
redis-cli -p 16386 SADD b2 c3 a1 
redis-cli -p 16386 SADD c3 a1 b2 
