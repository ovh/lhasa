all: dep
	ng build --output-path=../dist/webui --prod --aot

dev: dep
	ng build --output-path=../dist/webui

dep:
	npm install

test: dep
	ng test
	ng e2e
	
run: all
	npm start

live: 
	npm start

clean:
	echo "no cleaning command available"