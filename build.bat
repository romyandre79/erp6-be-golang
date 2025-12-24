@echo off
REM ===========================
REM BUILD
REM ===========================
echo Build To All OS
powershell -ExecutionPolicy Bypass -File build-all.ps1 

REM ===========================
REM COPY PUBLIC BE
REM ===========================
echo Copy Public BE to dist Public
xcopy public\* ..\erp6-be-golang-dist\public\ /E /I /Y

REM ===========================
REM COPY CONFIG BE
REM ===========================
echo Copy Config BE to dist Config
xcopy config\* ..\erp6-be-golang-dist\config\ /E /I /Y