# surveillance-cameras
Scripts for embedded systems processing surveillance cameras' data

### cameraSort.go
Reads the content of the directory specified in the first argument (origin), takes all the files that starts with CMIDyyyyMMdd (e.g. CAM120180217...) and puts them in destiny(second argument)/CAMID/yyyy/MM/dd

### rmOldCamera.go
Read the content of the first argument (path), which is structured like [path/]CAMID/yyyy/MM/dd/ and removes all the folders that are N (second argument) days old

### listLargeDirs.go
Lists a large directory, not recursively, printing the result with the following format: FILENAME @ IsDir? [true-false]

### fetchImage.go
Every second requests the picture provided in the URL, login with the credentials (if provided) and saves it in PATH/CamName/YYYY/MM/DD/ as hh-mm-ss.extension

### fetchVideo.go
Gets the videos stored in the internal memory of a IP Camera OneWay OWIPCAM45. It requires the camera's address (e.g. http://192.168.1.2), and an user account and password. It will save the videos in PATH/CamName/YYYY/MM/DD/ as hhmmss_hhmmss.extension, being the first "hhmmss" the starting time of the video, and the second the ending time.

### sort_APPIP01WV4
Reads the files from the first argument (origin) and moves them to destination(second argument)/CamName/YYYY/MM/DD. Is only valid for the pictures taken by the IP Camera Approx APPIP01WV4 that have the following format: "MacAddress00(NAME)_0_YYYYMMDDhhmmss_number.jpg".
