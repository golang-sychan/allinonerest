SHELL := /bin/bash

#获取当前文件所在路径（相对路径）
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

#如果工程根目录没定义，则根据当前路径反推出来
#cd $(COMMON_SELF_DIR)/../.. && pwd -P 获取工程目录的相对路径
# $(shell cd ...) 执行 shell命令， cd不是makefile 指令，所有要用shell函数包装下
#
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif

#定义可执行文件输出路径
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/output
endif

# set the version number. you should not need to do this
# for the majority of scenarios.
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --dirty --always --tags | sed 's/-/./2' | sed 's/-/./2' )
endif
export VERSION