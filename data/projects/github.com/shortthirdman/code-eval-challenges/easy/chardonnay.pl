use strict;

sub contains {
    my ( $s, $p ) = @_;
    while ( length $s > 0 && length $p > 0 ) {
        return (0) if ( substr( $s, 0, 1 ) > substr( $p, 0, 1 ) );
        $p = substr( $p, 1 ) if ( substr( $s, 0, 1 ) eq substr( $p, 0, 1 ) );
        $s = substr( $s, 1 );
    }
    return ( $p eq '' );
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my @s = split( / \| /, $line );
    my @t = split( / /,    $s[0] );
    my $p = join( '', sort split( //, $s[1] ) );
    my @u = map { join( '', sort split( //, $_ ) ) } @t;
    my @r;
    for ( my $i = 0 ; $i < scalar @t ; $i++ ) {
        push @r, $t[$i] if contains( $u[$i], $p );
    }
    printf( "%s\n", scalar @r == 0 ? "False" : join( ' ', @r ) );
}
close(INFILE);
