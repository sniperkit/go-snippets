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
        @s[$i] = @s[$i] % 10 + 1 if @s[$i] > 9;
    }
    my $r;
    map { $r += $_ } @s;
    print( ( $r % 10 == 0 ) ? "1\n" : "0\n" );
}
close(INFILE);
