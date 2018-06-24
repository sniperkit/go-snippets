use strict;
use List::Util qw(first);
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @s = split( / /, $line );
    my %h;
    $h{$_}++ for (@s);
    my $r = 0;
    my @t = sort keys %h;
    for ( my $i = 0 ; $i <= $#t ; $i++ ) {
        if ( $h{ $t[$i] } == 1 ) {
            $r = 1 + first { $s[$_] eq $t[$i] } 0 .. $#s;
            last;
        }
    }
    print "$r\n";
}
close(INFILE);
