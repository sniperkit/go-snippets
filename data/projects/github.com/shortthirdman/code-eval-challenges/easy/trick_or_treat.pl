use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my ( $v, $z, $w, $h ) = ( $line =~
          /Vampires: (\d+), Zombies: (\d+), Witches: (\d+), Houses: (\d+)/ );
    my $r =
      ( $v + $z + $w == 0 )
      ? 0
      : ( $v * 3 + $z * 4 + $w * 5 ) * $h / ( $v + $z + $w );
    print "$r\n";
}
close(INFILE);
