:: make sure it runs from its own directory
@echo off
cd /d "%~dp0"
:: deletes previous nep from where go installs it, then builds main as nep, then install it
powershell -Command "del %USERPROFILE%\go\bin\nep.exe; go build -o nep.exe; go install"