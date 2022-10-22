module main

go 1.16

require (
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/dustin/go-humanize v1.0.0
	github.com/miekg/dns v1.1.50
	github.com/pschou/go-rpm v0.0.0-00010101000000-000000000000
	github.com/pschou/go-tease v0.0.0-20220501223706-350142e428cd
	github.com/ulikunitz/xz v0.5.10
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	pault.ag/go/debian v0.12.0
)

replace pault.ag/go/debian => /home/schou/git/go-debian

replace pault.ag/go/debian/deb => /home/schou/git/go-debian/deb

replace github.com/pschou/go-rpm => /home/schou/git/go-rpm

replace github.com/pschou/go-tease => /home/schou/git/go-tease
