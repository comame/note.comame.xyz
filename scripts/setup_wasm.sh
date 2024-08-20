#! /bin/bash

go install github.com/agnivade/wasmbrowsertest@latest

if !(type google-chrome); then
    sudo wget -P /tmp https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
    sudo apt install /tmp/google-chrome-stable_current_amd64.deb
fi
