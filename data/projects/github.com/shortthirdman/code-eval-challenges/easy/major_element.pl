use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @s = split( /,/, $line );
    my $n = scalar @s / 2;
    my %t;
    map { $t{$_}++; } @s;
    my $f = 0;
    foreach my $i ( keys %t ) {
        next if $t{$i} < $n;
        last if $t{$i} == $n;
        print "$i\n";
        $f = 1;
        last;
    }
    print "None\n" if !$f;
}
close(INFILE);
