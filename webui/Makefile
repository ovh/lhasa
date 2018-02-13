all: dep
	ng build --output-path=../dist/webui

dep:
	npm install

test: dep
	ng test
	ng e2e
	
run: all
	ng serve

clean:
	echo "no cleaning command available"