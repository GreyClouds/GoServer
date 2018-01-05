@echo off
set base_dir=%~dp0
cd /d %~dp0
protoc --csharp_out=./hero/ hero.proto
protoc --go_out=./hero/ hero.proto
rem protoc --go_out=./go/ *.proto 

pause