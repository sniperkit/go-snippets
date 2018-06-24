use strict;
for ( my $i = 1 ; $i <= 12 ; $i++ ) {
    for ( my $j = 1 ; $j <= 12 ; $j++ ) {
        printf( "%4d", $i * $j );
    }
    print "\n";
}
