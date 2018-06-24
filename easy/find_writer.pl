use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my ( $s, $t ) = split( /\| /, $line );
    my @r = map { substr( $s, $_ - 1, 1 ) } split( / /, $t );
    printf( "%s\n", join( '', @r ) );
}
close(INFILE);
