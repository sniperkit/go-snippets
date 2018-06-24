use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    $line =~ s/^[0-9]+;//g;
    my @a;
    foreach ( split( /,/, $line ) ) {
        if ( @a[$_] ) {
            print "$_\n";
            last;
        }
        @a[$_] = 1;
    }
}
close(INFILE);
