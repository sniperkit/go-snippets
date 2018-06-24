use strict;

sub incircle {
    my ( $a, $b, $c ) = @_;
    return $a * $a + $b * $b <= $c * $c;
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    $line =~ tr/(),;//d;
    my @s = split( / /, $line );
    printf( incircle( $s[1] - $s[6], $s[2] - $s[7], $s[4] )
        ? "true\n"
        : "false\n" );
}
close(INFILE);
