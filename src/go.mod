module main

go 1.16

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/cavaliergopher/rpm v1.2.0
	github.com/dustin/go-humanize v1.0.0
	github.com/miekg/dns v1.1.50
	github.com/pschou/go_tease v0.0.0-20220501223706-350142e428cd
	github.com/ulikunitz/xz v0.5.10
	pault.ag/go/debian v0.12.0
)

replace pault.ag/go/debian => /home/schou/git/go-debian

replace pault.ag/go/debian/deb => /home/schou/git/go-debian/deb
