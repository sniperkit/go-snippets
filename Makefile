.PHONY: publish
all:
	@echo "type make publish"

update:
	git pull origin master
view:
	cd docs
	python -m SimpleHTTPServer 8000

publish:
	git commit -a -m "Added links"
	git push origin master
