all:
	make -C api
	make -C webui

test: 
	make -C api test integration-test
	make -C webui test
	
run: all
	./dist/appcatalog

clean:
	make -C api clean
	make -C webui clean
