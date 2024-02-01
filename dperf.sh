for i in 4 16 64 256 4096 65536; do
	for j in 1 4 16; do
		echo " dperf -b $i KiB -f 1GiB -i $j /borkdisk{1..12} --verbose  >> dperf-block_kb$i-threads-$j.out; "
    	do echo " dperf -b $i KiB -f 1GiB -i $j /borkdisk{1..12} --verbose --serial >> dperf-block_kb$i-threads-$j-serial.out; "
        do echo "i=$i and j=$j and k=$k"
                done;
                done;
        done;