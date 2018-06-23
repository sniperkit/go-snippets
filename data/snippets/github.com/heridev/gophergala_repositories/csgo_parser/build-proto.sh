#!/bin/bash

protoc --proto_path=proto --proto_path=/usr/include --go_out=. proto/cstrike15_usermessages_public.proto
