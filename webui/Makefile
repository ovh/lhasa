all: ui-kit
	ng build --output-path=../dist/webui --prod --aot --environment prod --base-href '{{ .UIBasePath }}' --deploy-url ${UI_BASE_HREF}

dev: ui-kit dep
	ng build --output-path=../dist/webui

ui-kit: dep
	@rm -rf ../dist/ovh-ui-kit
	@mkdir -p ../dist/ovh-ui-kit
	@echo copy ovh-ui-kit dist icons
	@cp -rfp ./node_modules//ovh-ui-kit/dist/icons ../dist/ovh-ui-kit
	@echo copy ovh-ui-kit dist fonts
	@cp -rfp ./node_modules//ovh-ui-kit/packages/oui-typography/fonts ../dist/ovh-ui-kit
	@echo less compile
	@./node_modules/@angular/cli/node_modules/.bin/lessc ./node_modules/ovh-ui-kit/packages/oui/stylekit.less -rp=./ --remcalc > ../dist/ovh-ui-kit/ovh-ui-kit.css
	@echo trick relative resouce location
	@sed s:../../dist/::g ../dist/ovh-ui-kit/ovh-ui-kit.css > ../dist/ovh-ui-kit/ovh-ui-kit.css.tmp
	@mv -f ../dist/ovh-ui-kit/ovh-ui-kit.css.tmp ../dist/ovh-ui-kit/ovh-ui-kit.css

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
