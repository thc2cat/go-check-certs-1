#include ../make/Makefile-for-go.mk

TAG=$(shell git tag)
NAME= $(notdir $(shell pwd))
SOURCES= $(wildcard *.go)
OPTIONS=-days 45

CLEANLIST=${NAME} ${NAME}-${TAG} ${NAME}.exe

all:	${NAME}

${NAME}:  ${SOURCES}
	go build

clean:
	@touch ${CLEANLIST}
	@rm ${CLEANLIST}
	

link:	${NAME}	
	ln -s ${NAME} ${NAME}-${TAG}

test: Services.txt  ${NAME} Makefile
	 ./${NAME} ${OPTIONS} 
