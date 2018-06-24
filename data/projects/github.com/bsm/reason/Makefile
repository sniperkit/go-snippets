PKG=$(subst github.com/bsm/reason,.,$(shell go list ./... | grep -v 'vendor'))
PROTOS=$(shell find .  -name '*.proto' -not -path './vendor/*')

default: vet test

vet:
	go vet $(PKG)

test:
	go test $(PKG)

test-short:
	go test $(PKG) -short

bench:
	go test $(PKG) -test.run=NONE -test.bench=. -benchmem -benchtime=1s

bench-race:
	go test $(PKG) -test.run=NONE -test.bench=. -race

.PHONY: vet test test-short bench bench-race

proto: proto.go proto.java
proto.go: $(patsubst %.proto,%.pb.go,$(PROTOS))
proto.java: java/lib/reason-java.jar

.PHONY: proto proto.go proto.java

# ---------------------------------------------------------------------

PROTO_PATH=.:vendor:vendor/github.com/gogo/protobuf/protobuf:../../..

### proto.go

%.pb.go: %.proto
	protoc --gogo_out=. --proto_path=${PROTO_PATH} $<

### proto.java

PROTO_JAVA_VERSION=3.5.1
noop=
space = $(noop) $(noop)

java/lib/reason-java.jar: java/dst/com/blacksquaremedia/reason/CoreProtos.class
	@mkdir -p $(dir $@)
	jar -cf $@ -C java/dst .

java/dst/com/blacksquaremedia/reason/CoreProtos.class: \
		java/lib/protobuf-java-$(PROTO_JAVA_VERSION).jar \
		java/lib/lz4-1.3.0.jar \
		java/src/com/blacksquaremedia/reason/CoreProtos.java \
		java/src/com/blacksquaremedia/reason/UtilProtos.java \
		java/src/com/blacksquaremedia/reason/classification/Hoeffding.java \
		java/src/com/blacksquaremedia/reason/classification/HoeffdingProtos.java \
		java/src/com/blacksquaremedia/reason/classification/FTRL.java \
		java/src/com/blacksquaremedia/reason/classification/FTRLProtos.java \
		java/src/com/blacksquaremedia/reason/regression/Hoeffding.java \
		java/src/com/blacksquaremedia/reason/regression/HoeffdingProtos.java \
		java/src/com/google/protobuf/GoGoProtos.java \
		$(shell find java/src -name '*.java')
	@mkdir -p $(dir $@)
	javac -cp $(subst $(space),:,$(filter %.jar,$^)) -d java/dst/ $(filter %.java,$^)

java/lib/protobuf-java-$(PROTO_JAVA_VERSION).jar:
	@mkdir -p  $(dir $@)
	curl -sSL https://repo1.maven.org/maven2/com/google/protobuf/protobuf-java/$(PROTO_JAVA_VERSION)/protobuf-java-$(PROTO_JAVA_VERSION).jar > $@

java/lib/lz4-1.3.0.jar:
	@mkdir -p  $(dir $@)
	curl -sSL https://repo1.maven.org/maven2/net/jpountz/lz4/lz4/1.3.0/lz4-1.3.0.jar > $@

java/src/com/blacksquaremedia/reason/CoreProtos.java: core/core.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<

java/src/com/blacksquaremedia/reason/UtilProtos.java: util/util.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<

java/src/com/blacksquaremedia/reason/classification/HoeffdingProtos.java: classification/hoeffding/internal/internal.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<

java/src/com/blacksquaremedia/reason/classification/FTRLProtos.java: classification/ftrl/internal/internal.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<

java/src/com/blacksquaremedia/reason/regression/HoeffdingProtos.java: regression/hoeffding/internal/internal.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<

java/src/com/google/protobuf/GoGoProtos.java: vendor/github.com/gogo/protobuf/gogoproto/gogo.proto
	@mkdir -p $(dir $@)
	protoc --java_out=java/src --proto_path=$(PROTO_PATH) $<
