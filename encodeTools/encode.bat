echo %1
echo %2
ffmpeg -f concat -safe 0 -r 10 -i %1 -vsync vfr -pix_fmt yuv420p %2