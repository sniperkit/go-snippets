use strict;

sub rotate {
    substr( @_[0], 1 ) . substr( @_[0], 0, 1 );
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my @s = split( /,/, $line );
    my $r = "False";
    for ( my $i = 0 ; $i < length $s[0] ; $i++ ) {
        if ( $s[0] eq $s[1] ) {
            $r = "True";
            last;
        }
        $s[1] = rotate( $s[1] );
    }
    print "$r\n";
}
close(INFILE);
