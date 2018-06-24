use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @s = split( / /, $line );
    my $n = shift @s;
    @s = sort { $a <=> $b } @s;
    my $u = @s[ $n >> 1 ];
    my $r = 0;
    map { $r += abs( $_ - $u ) } @s;
    print "$r\n";
}
close(INFILE);
