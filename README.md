# surveillance-cameras
Scripts for embedded systems processing surveillance cameras' data

### cameraSort.go
Reads the content of the directory specified in the first argument (origin), takes all the files that starts with CAMIDyyyyMMdd (e.g. Interior20180217...) and puts them in destiny(second argument)/CAMID/yyyy/MM/dd

### rmOldCamera.go
Read the content of the first argument (path), which is structured like [path/]CAMID/yyyy/MM/dd/ and removes all the folders that are N (second argument) days old

### listLargeDirs.go
Lists a large directory, not recursively, printing the result with the following format: FILENAME @ IsDir? [true-false]

### fetchImage.go
Every second requests the picture provided in the URL, login with the credentials (if provided) and saves it in PATH/CamName/YYYY/MM/DD/ as hh-mm-ss.extension
