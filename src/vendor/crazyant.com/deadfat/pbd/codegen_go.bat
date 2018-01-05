@echo off
set base_dir=%~dp0
cd /d %~dp0
protoc --go_out=./go/ id.proto fat.proto
rem protoc --go_out=./go/ fat.proto
rem protoc --go_out=./go/ *.proto 

pause