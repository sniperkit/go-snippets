use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my @a = ( $line =~ /(\d+):(\d+):(\d+) (\d+):(\d+):(\d+)/ );
    my $t =
      abs( ( @a[0] - @a[3] ) * 3600 + ( @a[1] - @a[4] ) * 60 + @a[2] - @a[5] );
    printf( "%02d:%02d:%02d\n", $t / 3600, ( $t / 60 ) % 60, $t % 60 );
}
close(INFILE);
