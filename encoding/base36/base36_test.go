package base36

import (
  "testing"
)

TestEncodeEmptyString( t *testing.T ) {
    should.ok(base62.encode('') === "", "OK");
  }

TestEncodeNil( t *testing.T ) {
    should.ok(base62.encode(undefined) === "", "OK");
  }

TestEncode_test123() {
    should.ok(base62.encode('test123') === "T6LpT34oC3", "OK");
  }

TestDecode_test123() {
    should.ok(base62.decode('T6LpT34oC3') === "test123", "OK");
  }

TestEncodeDecodeArray() {
  var bytes1 = [0x53, 0xFE, 0x92];
  var s1 = base62.encode(bytes1);

  it("encode arrray to string", function() {
    should.ok(s1 === 'Kzya2', "OK");
  });

  var bytes = [116, 32, 8, 99, 100, 232, 4, 7];

  // T208OsJe107
  var s = base62.encode(bytes);

  it("encode arrray to string", function() {
    should.ok(s === 'T208OsJe107', "OK");
  });

  var arr = base62.decodeToBuffer(s);

  it("decode string to array", function() {
    should.ok(new Buffer(bytes).toString() === arr.toString(), "OK");
  });
}

TestEncodeDecodeStr256() {
  var str256 = '';
  for (var i = 0; i <= 255; i++) {
    str256 += String.fromCharCode(i);
  }

  var strB62 = base62.encode(str256);
  //console.log(strB62);

  var strData = base62.decode(strB62);
  //console.log(strData);

  var boolReturn = (strData === str256);

  it("encode/decode 0-255", function() {
    should.ok(boolReturn, "OK");
  });
}

TestKeyBuf() {
  // this buff casued an encoding overrun beciase
  // last entry has zero length
  var keybuf = [24, 23, 224, 166, 164, 198, 162, 13, 94, 181, 12, 245, 108,
    24, 143, 220, 152, 181, 9, 74, 70, 81, 227, 157, 1, 41, 78, 125, 143,
    229, 88, 105, 247, 107, 128, 90, 144, 179, 55, 168, 51, 205, 190, 33,
    46, 123, 86, 123, 129, 206, 185, 206, 231, 48, 21, 76
  ]

  it("encode/decode keybuf example", function() {
    let key = base62.encode(keybuf);
    let decoded = base62.decodeToBuffer(key)
    should(decoded.length).be.equal(keybuf.length);
    should(decoded.equals(new Buffer(keybuf))).be.true;
  });
}

TestHighByteVals() {

    let val63 = [252]
    let encoded = base62.encode(val63);
    let decoded = base62.decodeToBuffer(encoded)
    should(decoded.equals(new Buffer(val63))).be.true;

    var val62 = [248]

    encoded = base62.encode(val62);
    decoded = base62.decodeToBuffer(encoded)
    should(decoded.equals(new Buffer(val62))).be.true;

    var val61 = [244]
    encoded = base62.encode(val61);
    decoded = base62.decodeToBuffer(encoded)
    should(decoded.equals(new Buffer(val61))).be.true;
  }
