V1 := data/heath/book01.json data/heath/book02.json
V2 := data/heath/book03.json data/heath/book04.json data/heath/book05.json data/heath/book06.json data/heath/book07.json data/heath/book08.json data/heath/book09.json
V3 := data/heath/book10.json data/heath/book11.json data/heath/book12.json data/heath/book13.json

all: $(V1) $(V2) $(V3)

# Intermediate trick from here https://stackoverflow.com/a/10609434
$(V1): intermediate.v1 ;
$(V2): intermediate.v2 ;
$(V3): intermediate.v3 ;
.INTERMEDIATE: intermediate.v1 intermediate.v2 intermediate.v3

intermediate.v1: heath/heath heath/vol1.xml
	./heath/heath -d data/heath heath/vol1.xml

intermediate.v2: heath/heath heath/vol2.xml
	./heath/heath -d data/heath heath/vol2.xml

intermediate.v3: heath/heath heath/vol3.xml
	./heath/heath -d data/heath heath/vol3.xml

heath/heath: heath/main.go heath/dom.go
	go build -o heath/heath ./heath

.PHONY: clean
clean:
	rm -rf data/heath
