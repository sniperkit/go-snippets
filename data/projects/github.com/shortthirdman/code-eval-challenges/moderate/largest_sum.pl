use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my @s = split( /,/, $line );
    my ( $l, $m ) = ( 0, @s[0] );
    for (@s) {
        $m = $_ + $l if $_ + $l > $m;
        $l = $_ + $l > 0 ? $_ + $l : 0;
    }
    print "$m\n";
}
close(INFILE);
