use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my @s = split( / /, $line );
    my $r;
    map {
        if ( length $_ > length $r ) { $r = $_ }
    } @s;
    print "$r\n";
}
close(INFILE);
