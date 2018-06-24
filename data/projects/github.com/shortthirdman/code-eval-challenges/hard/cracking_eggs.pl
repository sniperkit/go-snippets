use strict;

sub floors {
    my ( $e, $s ) = @_;
    return (1)  if $s == 1;
    return ($s) if $e == 1;
    return (0)  if $e == 0 || $s == 0;
    return floors( $s, $s ) if $e > $s;
    return floors( $e - 1, $s - 1 ) + floors( $e, $s - 1 ) + 1;
}

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $e, $f ) = split( / /, $line );
    my $n = 0;
    $n++ while floors( $e, $n ) < $f;
    print "$n\n";
}
close(INFILE);
