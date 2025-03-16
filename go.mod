module cloudartisan.com/lexo

go 1.20

require github.com/abadojack/whatlanggo v1.0.1

// Force correct versions
replace (
	golang.org/x/text => golang.org/x/text v0.3.8
	github.com/boyter/scc => github.com/boyter/scc v1.12.1
)
