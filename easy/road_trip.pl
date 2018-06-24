use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    $line =~ s/[A-Z][a-z]+,|;//g;
    my $s = 0;
    my @r = map {
        my $t = $_ - $s;
        $s += $t;
        $t;
    } sort { $a <=> $b } split( / /, $line );
    printf( "%s\n", join( ',', @r ) );
}
close(INFILE);
