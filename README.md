# surveillance-cameras
Scripts for embedded systems processing surveillance cameras' data

### APPIP01WV4_sort
Sorts the pictures taken by a Approx APPIP01WV4.

### CNETCAM_sort
Sorts the pictures taken by a Conceptronic CNETCAM.

### generic_listLargeDirs
List (unsorted) very large directories.

### generic_recordVideo
Records video of a IP Camera using RTSP, and stores it in (path/)Name/YYYY/MM/DD as hh-mm-ss.mkv. It requires ffmpeg and the RTSP stream MUST have codecs that are MKV compatible.

### generic_rmOldCameraData
In a path structured like (path/)Name/YYYY/MM/DD, it removes the directories that are older than ~ days \[default=30\].

### OWIPCAM4X_fetchImage
Fetch a picture every second to a OneWay IP Camera (tested with models OWIPCAM43 and OWIPCAM45) and save it in (path/)CameraName/YYYY/MM/DD/hh-mm-ss.ext. 

### OWIPCAM45_fetchVideo
Fetch all videos stored in a OneWay OWIPCAM45's internal memory.

### OWIPCAM45_rotate
Rotates a OneWay OWIPCAM45 from left to right, and from right to left, with ~ rotations \[default=10\].
