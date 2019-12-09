#!/bin/bash

TESTS_PATH="${GOPATH}/src/github.com/Miguel-Dorta/surveillance-cameras/test"

"${TESTS_PATH}/APPIP01WV4_sort/test.sh"
"${TESTS_PATH}/CNETCAM_sort/test.sh"
"${TESTS_PATH}/generic_listLargeDirs/test.sh"
"${TESTS_PATH}/generic_rmOldCameraData/test.sh"
"${TESTS_PATH}/OWIPCAM4X_fetchImage/test.sh"
