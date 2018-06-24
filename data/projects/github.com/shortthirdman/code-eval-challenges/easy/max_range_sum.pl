use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $n, $s ) = split( /;/, $line );
    my @t = split( / /, $s );
    my ( $m, $c ) = ( 0, 0 );
    for ( my $i = 0 ; $i < $n ; $i++ ) {
        $c += @t[$i];
    }
    $m = $c if ( $c > $m );
    for ( my $i = $n ; $i < scalar @t ; $i++ ) {
        $c += @t[$i] - @t[ $i - $n ];
        $m = $c if ( $c > $m );
    }
    print "$m\n";
}
close(INFILE);
