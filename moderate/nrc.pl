use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my %seen;
    $seen{$_}++ foreach split( //, $line );
    my $r;
    foreach my $i ( split( //, $line ) ) {
        if ( $seen{$i} == 1 ) { $r = $i; last }
    }
    print "$r\n";
}
close(INFILE);
