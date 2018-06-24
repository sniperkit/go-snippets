use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    $line =~ s/ //g;
    my @s = split( //, $line );
    for ( my $i = scalar @s - 2 ; $i >= 0 ; $i -= 2 ) {
        @s[$i] <<= 1;
    }
    my $r;
    map { $r += $_ } @s;
    print( ( $r % 10 == 0 ) ? "Real\n" : "Fake\n" );
}
close(INFILE);
