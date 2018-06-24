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
    my @a = map { "*" x $_ . substr( $r, $_, 1 ) } 0 .. ( length $r ) - 1;
    printf( "%s\n", join( ' ', @a ) );
}
close(INFILE);
