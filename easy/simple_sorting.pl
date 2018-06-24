use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @r = sort { $a <=> $b } split( / /, $line );
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
