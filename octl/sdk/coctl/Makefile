all: test lib

test: lib 
	gcc -g -o test test.c libcoctl.a
	gcc -g -o test_shared test.c -lcoctl -L.
	
lib: wrapper.a coctl.o
	ar -x wrapper.a
	ar -r libcoctl.a *.o
ifeq ($(OS), Windows_NT)
	gcc -shared -o libcoctl.dll *.o	
else
	gcc -shared -o libcoctl.so *.o
endif

wrapper.a: wrapper.go
	go build -buildmode=c-archive -o wrapper.a wrapper.go

coctl.o: coctl.c coctl.h wrapper.a
	gcc -c -fPIC coctl.c

install: lib
ifeq ($(OS), Windows_NT)
	echo "target <install> not suitable for Windows"	
else
	mkdir /usr/local/include/coctl
	cp *.h /usr/local/include/coctl/
	cp libcoctl.so libcoctl.a /usr/local/lib/
	ldconfig
	echo "Installing done."
endif
	

uninstall:
ifeq ($(OS), Windows_NT)
	echo "target <uninstall> not suitable for Windows"	
else
	rm -rf /usr/local/include/coctl
	rm /usr/local/lib/libcoctl.so /usr/local/lib/libcoctl.a
	ldconfig
	echo "Uninstalling done."
endif

clean:
ifeq ($(OS), Windows_NT)
	del *.o
	del *.gch
	del *.a 
	del test 
	del test_shared 
	del *.dll
else
	rm -rf *.o *.gch *.a test test_shared *.so
endif

.PHONY: clean
