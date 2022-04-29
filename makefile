SHELL:=/bin/bash
PROJECT:=preppy

fmt:
	@go fmt

fetch.stdlib:
	@python3 stdlib.py

build: fmt fetch.stdlib
	@go build -ldflags="-s -w" && upx ${PROJECT}

test.setup:
	@git clone https://github.com/ashleve/lightning-hydra-template.git ashleve && rm -rf ashleve/.git
	@git clone https://github.com/fastai/nbdev.git nbdev && rm -rf nbdev/.git

test.ashleve: build
	@cd ashleve && time ../${PROJECT} -d

test.nbdev: build
	@time ./${PROJECT} -d -r nbdev