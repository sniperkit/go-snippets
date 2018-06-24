use strict;
use Config;
if ( $Config{byteorder} =~ m/^1/ ) {
    print "LittleEndian\n";
}
else {
    print "BigEndian\n";
}
