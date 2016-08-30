#!/bin/bash

(cd news && protoc --go_out=plugins=grpc:. *.proto)
(cd profile && protoc --go_out=plugins=grpc:. *.proto)
