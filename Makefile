
GOFMT=gofmt -s -spaces=true -tabindent=false -tabwidth=4

format:
	${GOFMT} -w oauth/oauth.go
	${GOFMT} -w golangtw/golangtw.go
	${GOFMT} -w util/util.go
	${GOFMT} -w util/sina.go



