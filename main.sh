while true
do
	./main &
	pid = $!
	sleep 15
	kill $pid
done