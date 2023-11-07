all: test lib

test: lib 
	gcc -g -o test test.c libcoctl.a
	gcc -g -o test_shared test.c -lcoctl -L.
	
lib: wrapper.a coctl.o
	ar -x wrapper.a
	ar -r libcoctl.a *.o
	gcc -shared -o libcoctl.so *.o

wrapper.a: wrapper.go
	go build -buildmode=c-archive -o wrapper.a wrapper.go
	go build -buildmode=c-shared -o libwrapper.so wrapper.go

coctl.o: coctl.c coctl.h
	gcc -c -fPIC coctl.c

install: lib
	mkdir /usr/local/include/coctl
	cp *.h /usr/local/include/coctl/
	cp libcoctl.so libcoctl.a /usr/local/lib/
	ldconfig
	echo "Installing done."

uninstall:
	rm -rf /usr/local/include/coctl
	rm /usr/local/lib/libcoctl.so /usr/local/lib/libcoctl.a
	ldconfig
	echo "Uninstalling done."
clean:
	rm -rf *.o *.gch *.a test test_shared *.so
