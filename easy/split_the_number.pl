use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my ( $l, $r ) = split( / /, $line );
    my ( $c, $s, $v, $o ) = ( 0, 0, 0, 1 );
    map {
        if ( $_ =~ /[a-z]/ )
        {
            $v = $v * 10 + substr( $l, $c, 1 );
            $c++;
        }
        elsif ( $_ eq '+' ) {
            $s += $v * $o;
            ( $v, $o ) = ( 0, 1 );
        }
        elsif ( $_ eq '-' ) {
            $s += $v * $o;
            ( $v, $o ) = ( 0, -1 );
        }
    } split( //, $r );
    printf( "%d\n", $s + $v * $o );
}
close(INFILE);
