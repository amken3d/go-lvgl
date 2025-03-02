BIN = golvgl
GOCMD = go
GOBUILD = ${GOCMD} build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"'
GOCLEAN = ${GOCMD} clean

build:
	${GOBUILD} -o ${BIN} -v *.go

clean:
	${GOCLEAN}
	rm -f ./lvgl/*_gen.go

init:
	git submodule init
	git submodule update
	mkdir -p include/obj/lvgl/lv_drivers/display
	mkdir -p include/obj/lvgl/lv_drivers/gtkdrv
	mkdir -p include/obj/lvgl/lv_drivers/indev
	ln -s -f ../lv_drivers include/lvgl

lvgl:
	$(MAKE) -C include all

clean-lvgl:
	$(MAKE) -C include clean

clean-all: clean clean-lvgl

generate:
	cd lvgl; ${GOCMD} generate

.PHONY: build clean init lvgl clean-lvgl clean-all generate