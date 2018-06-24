use strict;
use List::Util qw[max];
my @a;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my @s = split( / /, $line );
    @s[0]  += @a[0]  if scalar @a > 0;
    @s[-1] += @a[-1] if scalar @a > 1;
    for ( my $i = 1 ; $i < scalar @a ; $i++ ) {
        @s[$i] += max @a[ $i - 1 ], @a[$i];
    }
    @a = @s;
}
printf( "%d\n", max @a );
close(INFILE);
