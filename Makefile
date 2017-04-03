VERSION=0.1.0
NAME=boxter_$(VERSION)
OUT_DIR=bin/linux_386/boxter_$(VERSION)

all:$(OUT_DIR)/boxter
$(OUT_DIR)/boxter:main.go
	gox  \
		-output "bin/{{.Dir}}_$(VERSION)/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		-osarch "linux/386" github.com/FarmRadioHangar/boxter

tar:
	cd bin/ && tar -zcvf boxter_$(VERSION).tar.gz  boxter_$(VERSION)/