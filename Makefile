#定义全局使用变量
ROOT_PACKAGE=github.com/golang-sychan/allinonerest
VERSION_PACKAGE=xxx/version

include build/lib/common.mk
include build/lib/golang.mk

.PHONY: build
build:
	@$(MAKE) go.build
