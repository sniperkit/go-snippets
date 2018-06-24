use strict;
use integer;

sub happy {
    my ( $a, $r ) = @_[0], 0;
    while ( $a > 0 ) {
        my $b = $a % 10;
        $r += $b * $b;
        $a /= 10;
    }
    return $r;
}

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my %b = ( $line => 1 );
    for ( my $i = 0 ; $i < 127 && $line > 1 ; $i++ ) {
        $line = happy($line);
        $line == 0 if ( $b{$line} );
        $b{$line} = 1;
    }
    printf( "%d\n", $line == 1 );
}
close(INFILE);
