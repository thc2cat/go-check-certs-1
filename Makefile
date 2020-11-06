TAG=$(shell git tag)
NAME= $(notdir $(shell pwd))
SOURCES= $(wildcard *.go)
OPTIONS=-days 90

CLEANLIST=${NAME} ${NAME}-${TAG} ${NAME}.exe

all:	${NAME}  results.txt

${NAME}:  ${SOURCES}
	go build

clean:
	@touch ${CLEANLIST}
	@rm ${CLEANLIST}
	

link:	${NAME}	
	ln -s ${NAME} ${NAME}-${TAG}

check: Services.txt  ${NAME} Makefile
	 @./${NAME} ${OPTIONS} 
