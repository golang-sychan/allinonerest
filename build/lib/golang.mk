C ?= cobra
ARCH ?= amd64
GO = go
#指定支持的golang版本
GO_SUPPORTED_VERSIONS ?= 1.11|1.12|1.13|1.14|1.15|1.16|1.18

GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) -X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

ifeq ($(OUTPUT_DIR),)
$(error the variable OUTPUT_DIR must be set prior to including golang.mk)
endif

ifeq ($(origin PLATFORM), undefined)
ifeq ($(origin OS), undefined)
OS := $(shell go env GOOS)
endif
ifeq ($(origin ARCH), undefined)
ARCH := $(shell go env GOARCH)
endif
PLATFORM := $(OS)_$(ARCH)
else
OS := $(word 1, $(subst _, ,$(PLATFORM)))
ARCH := $(word 2, $(subst _, ,$(PLATFORM)))
endif


#判断操作系统
ifeq ($(OS),windows)
GO_OUT_EXT := .exe
endif

#获取宿主机的OS，ARCH
GOHOSTOS := $(shell go env GOHOSTOS)
GOHOSTARCH := $(shell go env GOHOSTARCH)
HOST_PLATFORM := $(GOHOSTOS)_$(GOHOSTARCH)

#指定路径下的所有文件和文件夹
COMMANDS=$(wildcard ${ROOT_DIR}/cmd/*)
BINS=$(foreach cmd,${COMMANDS},$(notdir ${cmd}))


#确认golang的版本支持
#grep -q 静默grep标准输出 -E 指定正则表达式
# && echo 0 || echo 1 正常返回0， 异常返回 1
.PHONY: go.build.verify
go.build.verify:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif


# $* 表示%所匹配的部分
# 如 go.build.linux_arm64.user, 则这里的 $* 就是 linux_arm64.user
# $(subst ., ,$*) 为 linux_arm64 user
# $(eval COMMAND := $(word 2,$(subst ., ,$*))) COMMAND 为 user
# $(eval PLATFORM := $(word 1,$(subst ., ,$*))) PLATFORM 为 linux_arm64
# $(eval OS := $(word 1,$(subst _, ,$(PLATFORM)))) OS 为 linux
# $(eval OS := $(word 2,$(subst _, ,$(PLATFORM)))) OS 为 arm64

go.build.%:
	@echo $*
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/$(OS)/$(ARCH)
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) -ldflags "$(GO_LDFLAGS)" $(ROOT_DIR)/cmd/$(COMMAND)


#编译之前确认golang版本是否支持
#BINS为cmd下的go files
#go.build.linux_arm64.user
.PHONY: go.build
go.build: go.build.verify $(addprefix go.build., $(addprefix $(HOST_PLATFORM)., $(BINS)))
