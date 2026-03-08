# require apt-get install xvfb
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0

go test -v -covermode=atomic ./... 