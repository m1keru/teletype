# tool marcros
CC := go
CCFLAG := 

# path marcros
BIN_PATH := dist
OBJ_PATH := obj
SRC_PATH := cmd

# compile marcros
TARGET_NAME := teletype
ifeq ($(OS),Windows_NT)
	TARGET_NAME := $(addsuffix .exe,$(TARGET_NAME))
endif
TARGET := $(BIN_PATH)/$(TARGET_NAME)
MAIN_SRC := cmd/teletype/main.go

# src files & obj files
SRC := $(foreach x, $(SRC_PATH), $(wildcard $(addprefix $(x)/*,.c*)))
OBJ := $(addprefix $(OBJ_PATH)/, $(addsuffix .o, $(notdir $(basename $(SRC)))))

# clean files list
DISTCLEAN_LIST := $(OBJ)
CLEAN_LIST := $(TARGET) \
			  $(DISTCLEAN_LIST)

# default rule
default: all

# non-phony targets
$(TARGET): $(OBJ)
	cd cmd/$(TARGET_NAME) && \
	$(CC) build -o ../../dist/$(TARGET_NAME) && \
	cd - && \
	cp config/config.yaml.tpl dist/ && \
	cp build/install.sh dist/ && \
	cp init/teletype.service dist/ && \
	tar -zcf dist.tgz dist && \
	mv dist.tgz dist
	

# phony rules
.PHONY: all
all: $(TARGET)

.PHONY: clean
clean:
	@echo CLEAN $(CLEAN_LIST)
	@rm -f $(CLEAN_LIST)

.PHONY: distclean
distclean:
	@echo CLEAN $(CLEAN_LIST)
	@rm -f $(DISTCLEAN_LIST)

