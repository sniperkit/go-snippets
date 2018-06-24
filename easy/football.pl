use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my $i = 1;
    my @r;
    my @u;
    my %t;
    map {
        if ( $_ eq '|' ) { $i++; }
        else {
            if ( exists $t{$_} ) { $t{$_} .= ',' . $i; }
            else {
                $t{$_} = $i;
                push @u, $_;
            }
        }
    } split( / /, $line );
    @u = sort { $a <=> $b } @u;
    for my $j ( 0 .. $#u ) {
        push @r, sprintf( "%d:%s;", $u[$j], $t{ $u[$j] } );
    }
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
