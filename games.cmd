@echo off
setlocal

rem -----------------------------------------
rem Positional arguments
rem -----------------------------------------
set APP=%1

if "%APP%"=="" (
    echo Usage: make.cmd ^<app^|all>^>
    exit /b 1
)

set DIST=dist

rem -----------------------------------------
rem Build ALL apps
rem -----------------------------------------
if "%APP%"=="all" goto :build_all

rem -----------------------------------------
rem Build ONE app
rem -----------------------------------------
call :package_app %APP%
goto :eof


:build_all
for /D %%A in (cmd\*) do (
    call :package_app %%~nA
)
goto :eof


:package_app
set APP=%~1

echo.
echo === Packaging %APP% for Windows/%ARCH% ===

pushd cmd\%APP%
fyne package --release --icon=splash.png --id=com.scsi.%APP%
popd

if not exist %DIST% mkdir %DIST%
copy cmd\%APP%\%APP%.exe %DIST%\

echo Output: %DIST%\%APP%.exe
goto :eof
